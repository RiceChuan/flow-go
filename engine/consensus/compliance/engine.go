package compliance

import (
	"fmt"
	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/consensus/hotstuff/model"
	"github.com/onflow/flow-go/engine"
	"github.com/onflow/flow-go/engine/common/fifoqueue"
	"github.com/onflow/flow-go/engine/consensus"
	"github.com/onflow/flow-go/engine/consensus/sealing/counters"
	"github.com/onflow/flow-go/model/events"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/messages"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/module/component"
	"github.com/onflow/flow-go/module/irrecoverable"
	"github.com/onflow/flow-go/module/metrics"
	"github.com/onflow/flow-go/module/util"
	"github.com/onflow/flow-go/network"
	"github.com/onflow/flow-go/network/channels"
	"github.com/onflow/flow-go/state/protocol"
	"github.com/onflow/flow-go/storage"
)

// defaultBlockQueueCapacity maximum capacity of inbound queue for `messages.BlockProposal`s
const defaultBlockQueueCapacity = 10_000

// Engine is a wrapper around `compliance.Core`. The Engine queues inbound messages, relevant
// node-internal notifications, and manages the worker routines processing the inbound events,
// and forwards outbound messages to the networking layer.
// `compliance.Core` implements the actual compliance logic.
type Engine struct {
	log            zerolog.Logger
	mempoolMetrics module.MempoolMetrics
	engineMetrics  module.EngineMetrics
	me             module.Local
	headers        storage.Headers
	payloads       storage.Payloads
	tracer         module.Tracer
	state          protocol.State
	prov           consensus.ProposalProvider
	core           *Core
	// queues for inbound messsages
	pendingBlocks  engine.MessageStore
	messageHandler *engine.MessageHandler
	// tracking finalized view
	finalizedView              counters.StrictMonotonousCounter
	finalizationEventsNotifier engine.Notifier

	cm *component.ComponentManager
	component.Component
}

var _ network.MessageProcessor = (*Engine)(nil)
var _ component.Component = (*Engine)(nil)

func NewEngine(
	log zerolog.Logger,
	me module.Local,
	prov consensus.ProposalProvider,
	core *Core,
) (*Engine, error) {

	// Inbound FIFO queue for `messages.BlockProposal`s
	blocksQueue, err := fifoqueue.NewFifoQueue(
		fifoqueue.WithCapacity(defaultBlockQueueCapacity),
		fifoqueue.WithLengthObserver(func(len int) { core.mempoolMetrics.MempoolEntries(metrics.ResourceBlockProposalQueue, uint(len)) }),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create queue for inbound block proposals: %w", err)
	}
	pendingBlocks := &engine.FifoMessageStore{FifoQueue: blocksQueue}

	// define message queueing behaviour
	handler := engine.NewMessageHandler(
		log.With().Str("compliance", "engine").Logger(),
		engine.NewNotifier(),
		engine.Pattern{
			Match: func(msg *engine.Message) bool {
				_, ok := msg.Payload.(*messages.BlockProposal)
				if ok {
					core.engineMetrics.MessageReceived(metrics.EngineCompliance, metrics.MessageBlockProposal)
				}
				return ok
			},
			Store: pendingBlocks,
		},
		engine.Pattern{
			Match: func(msg *engine.Message) bool {
				_, ok := msg.Payload.(*events.SyncedBlock)
				if ok {
					core.engineMetrics.MessageReceived(metrics.EngineCompliance, metrics.MessageSyncedBlock)
				}
				return ok
			},
			Map: func(msg *engine.Message) (*engine.Message, bool) {
				syncedBlock := msg.Payload.(*events.SyncedBlock)
				msg = &engine.Message{
					OriginID: msg.OriginID,
					Payload: &messages.BlockProposal{
						Payload: syncedBlock.Block.Payload,
						Header:  syncedBlock.Block.Header,
					},
				}
				return msg, true
			},
			Store: pendingBlocks,
		},
	)

	eng := &Engine{
		log:                        log.With().Str("compliance", "engine").Logger(),
		me:                         me,
		mempoolMetrics:             core.mempoolMetrics,
		engineMetrics:              core.engineMetrics,
		headers:                    core.headers,
		payloads:                   core.payloads,
		pendingBlocks:              pendingBlocks,
		state:                      core.state,
		tracer:                     core.tracer,
		prov:                       prov,
		core:                       core,
		messageHandler:             handler,
		finalizationEventsNotifier: engine.NewNotifier(),
	}

	// create the component manager and worker threads
	eng.cm = component.NewComponentManagerBuilder().
		AddWorker(eng.processMessagesLoop).
		AddWorker(eng.finalizationProcessingLoop).
		Build()
	eng.Component = eng.cm

	return eng, nil
}

// WithConsensus adds the consensus algorithm to the engine. This must be
// called before the engine can start.
func (e *Engine) WithConsensus(hot module.HotStuff) *Engine {
	e.core.hotstuff = hot
	return e
}

// Start starts the Hotstuff event processMessagesLoop, then the compliance engine worker threads.
func (e *Engine) Start(ctx irrecoverable.SignalerContext) {
	if e.core.hotstuff == nil {
		ctx.Throw(fmt.Errorf("must initialize compliance engine with hotstuff engine"))
	}

	e.log.Info().Msg("starting hotstuff")
	e.core.hotstuff.Start(ctx)
	e.log.Info().Msg("hotstuff started")

	e.log.Info().Msg("starting compliance engine")
	e.Component.Start(ctx)
	e.log.Info().Msg("compliance engine started")
}

// Ready returns a ready channel that is closed once the engine has fully started.
// For the consensus engine, we wait for hotstuff to start.
func (e *Engine) Ready() <-chan struct{} {
	// NOTE: this will create long-lived goroutines each time Ready is called
	// Since Ready is called infrequently, that is OK. If the call frequency changes, change this code.
	return util.AllReady(e.cm, e.core.hotstuff)
}

// Done returns a done channel that is closed once the engine has fully stopped.
// For the consensus engine, we wait for hotstuff to finish.
func (e *Engine) Done() <-chan struct{} {
	// NOTE: this will create long-lived goroutines each time Done is called
	// Since Done is called infrequently, that is OK. If the call frequency changes, change this code.
	return util.AllDone(e.cm, e.core.hotstuff)
}

// Process processes the given event from the node with the given origin ID in
// a blocking manner. It returns the potential processing error when done.
func (e *Engine) Process(channel channels.Channel, originID flow.Identifier, event interface{}) error {
	err := e.messageHandler.Process(originID, event)
	if err != nil {
		if engine.IsIncompatibleInputTypeError(err) {
			e.log.Warn().Msgf("%v delivered unsupported message %T through %v", originID, event, channel)
			return nil
		}
		return fmt.Errorf("unexpected error while processing engine message: %w", err)
	}
	return nil
}

// processMessagesLoop processes available block, vote, and timeout messages as they are queued.
func (e *Engine) processMessagesLoop(ctx irrecoverable.SignalerContext, ready component.ReadyFunc) {
	ready()

	doneSignal := ctx.Done()
	newMessageSignal := e.messageHandler.GetNotifier()

	for {
		select {
		case <-doneSignal:
			return
		case <-newMessageSignal:
			err := e.processAvailableMessages() // no errors expected during normal operations
			if err != nil {
				ctx.Throw(err)
			}
		}
	}
}

// processAvailableMessages processes any available messages until the message queue is empty.
// Only returns when all inbound queues are empty (or the engine is terminated).
// No errors are expected during normal operation. All returned exceptions are potential
// symptoms of internal state corruption and should be fatal.
func (e *Engine) processAvailableMessages() error {
	for {
		msg, ok := e.pendingBlocks.Get()
		if ok {
			err := e.core.OnBlockProposal(msg.OriginID, msg.Payload.(*messages.BlockProposal))
			if err != nil {
				return fmt.Errorf("could not handle block proposal: %w", err)
			}
			continue
		}

		// when there are no more messages in the queue, back to the processMessagesLoop to wait
		// for the next incoming message to arrive.
		return nil
	}
}

// OnFinalizedBlock implements the `OnFinalizedBlock` callback from the `hotstuff.FinalizationConsumer`
// It informs sealing.Core about finalization of respective block.
//
// CAUTION: the input to this callback is treated as trusted; precautions should be taken that messages
// from external nodes cannot be considered as inputs to this function
func (e *Engine) OnFinalizedBlock(block *model.Block) {
	if e.finalizedView.Set(block.View) {
		e.finalizationEventsNotifier.Notify()
	}
}

// finalizationProcessingLoop is a separate goroutine that performs processing of finalization events
func (e *Engine) finalizationProcessingLoop(ctx irrecoverable.SignalerContext, ready component.ReadyFunc) {
	ready()

	doneSignal := ctx.Done()
	blockFinalizedSignal := e.finalizationEventsNotifier.Channel()
	for {
		select {
		case <-doneSignal:
			return
		case <-blockFinalizedSignal:
			e.core.ProcessFinalizedView(e.finalizedView.Value())
		}
	}
}
