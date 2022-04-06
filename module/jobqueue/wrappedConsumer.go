package jobqueue

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/module/component"
	"github.com/onflow/flow-go/module/irrecoverable"
	"github.com/onflow/flow-go/storage"
)

type JobPreProcessor func(interface{})

type WrappedConsumerOption func(*WrappedConsumer)

func WitNotifier(notifier NotifyDone) WrappedConsumerOption {
	return func(c *WrappedConsumer) {
		c.notifier = notifier
	}
}

func WithMaxProcessing(max uint64) WrappedConsumerOption {
	return func(c *WrappedConsumer) {
		c.maxProcessing = max
	}
}

type WrappedConsumer struct {
	component.Component
	module.Resumable

	cm           *component.ComponentManager
	consumer     module.JobConsumer
	jobs         module.Jobs
	workSignal   <-chan struct{}
	defaultIndex uint64
	log          zerolog.Logger

	preprocessor  JobPreProcessor
	notifier      NotifyDone
	maxProcessing uint64
}

// NewWrappedConsumer creates a new WrappedConsumer consumer
func NewWrappedConsumer(
	log zerolog.Logger,
	progress storage.ConsumerProgress,
	jobs module.Jobs,
	processor JobProcessor, // method used to process jobs
	workSignal <-chan struct{},
	defaultIndex uint64,
	opts ...WrappedConsumerOption,
) (*WrappedConsumer, error) {

	c := &WrappedConsumer{
		defaultIndex:  defaultIndex,
		workSignal:    workSignal,
		jobs:          jobs,
		log:           log,
		maxProcessing: 1,
	}

	for _, opt := range opts {
		opt(c)
	}

	worker := NewWrappedWorker(
		processor,
		func(id module.JobID) { c.NotifyJobIsDone(id) },
		c.maxProcessing,
	)
	c.consumer = NewConsumer(c.log, c.jobs, progress, worker, c.maxProcessing)

	builder := component.NewComponentManagerBuilder().
		AddWorker(func(ctx irrecoverable.SignalerContext, ready component.ReadyFunc) {
			worker.Start(ctx)
			<-worker.Ready()

			err := c.consumer.Start(c.defaultIndex)
			if err != nil {
				ctx.Throw(fmt.Errorf("could not start consumer: %w", err))
			}

			ready()

			<-ctx.Done()

			// blocks until all running jobs have stopped
			c.consumer.Stop()

			<-worker.Done()
		}).
		AddWorker(func(ctx irrecoverable.SignalerContext, ready component.ReadyFunc) {
			ready()
			c.processingLoop(ctx)
		})

	cm := builder.Build()
	c.cm = cm
	c.Component = cm
	c.Resumable = c.consumer

	return c, nil
}

// NotifyJobIsDone is invoked by the worker to let the consumer know that it is done
// processing a (block) job.
func (c *WrappedConsumer) NotifyJobIsDone(jobID module.JobID) uint64 {
	// notify wrapped consumer that job is complete
	c.defaultIndex = c.consumer.NotifyJobIsDone(jobID)

	// notify instantiator that job is complete
	if c.notifier != nil {
		c.notifier(jobID)
	}

	return c.defaultIndex
}

// Size returns number of in-memory block jobs that block consumer is processing.
func (c *WrappedConsumer) Size() uint {
	return c.consumer.Size()
}

// Head returns the highest job index available
func (c *WrappedConsumer) Head() (uint64, error) {
	return c.jobs.Head()
}

func (c *WrappedConsumer) processingLoop(ctx irrecoverable.SignalerContext) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.workSignal:
			c.consumer.Check()
		}
	}
}
