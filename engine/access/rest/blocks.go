package rest

import (
	"context"
	"fmt"

	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/onflow/flow-go/access"
	"github.com/onflow/flow-go/engine/access/rest/models"
	"github.com/onflow/flow-go/engine/access/rest/request"
	"github.com/onflow/flow-go/engine/common/rpc/convert"
	"github.com/onflow/flow-go/model/flow"
	accessproto "github.com/onflow/flow/protobuf/go/flow/access"
	"github.com/onflow/flow/protobuf/go/flow/entities"
)

// GetBlocksByIDs gets blocks by provided ID or list of IDs.
func GetBlocksByIDs(r *request.Request, srv RestServerApi, link models.LinkGenerator) (interface{}, error) {
	req, err := r.GetBlockByIDsRequest()
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	return srv.GetBlocksByIDs(req, r.Context(), r.ExpandFields, link)
}

// GetBlocksByHeight gets blocks by height.
func GetBlocksByHeight(r *request.Request, srv RestServerApi, link models.LinkGenerator) (interface{}, error) {
	return srv.GetBlocksByHeight(r, link)
}

// GetBlockPayloadByID gets block payload by ID
func GetBlockPayloadByID(r *request.Request, srv RestServerApi, link models.LinkGenerator) (interface{}, error) {
	req, err := r.GetBlockPayloadRequest()
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	return srv.GetBlockPayloadByID(req, r.Context(), link)
}

func getBlock(option blockRequestOption, context context.Context, expandFields map[string]bool, backend access.API, link models.LinkGenerator) (*models.Block, error) {
	// lookup block
	blkProvider := NewBlockRequestProvider(backend, option)
	blk, blockStatus, err := blkProvider.getBlock(context)
	if err != nil {
		return nil, err
	}

	// lookup execution result
	// (even if not specified as expandable, since we need the execution result ID to generate its expandable link)
	var block models.Block
	executionResult, err := backend.GetExecutionResultForBlockID(context, blk.ID())
	if err != nil {
		// handle case where execution result is not yet available
		if se, ok := status.FromError(err); ok {
			if se.Code() == codes.NotFound {
				err := block.Build(blk, nil, link, blockStatus, expandFields)
				if err != nil {
					return nil, err
				}
				return &block, nil
			}
		}
		return nil, err
	}

	err = block.Build(blk, executionResult, link, blockStatus, expandFields)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func getForwarderBlock(option blockRequestOption, context context.Context, expandFields map[string]bool, upstream accessproto.AccessAPIClient, link models.LinkGenerator) (*models.Block, error) {
	// lookup block
	blkProvider := NewBlockForwarderProvider(upstream, option)
	blk, blockStatus, err := blkProvider.getBlock(context)
	if err != nil {
		return nil, err
	}

	// lookup execution result
	// (even if not specified as expandable, since we need the execution result ID to generate its expandable link)
	var block models.Block
	getExecutionResultForBlockIDRequest := &accessproto.GetExecutionResultForBlockIDRequest{
		BlockId: blk.Id,
	}

	executionResultForBlockIDResponse, err := upstream.GetExecutionResultForBlockID(context, getExecutionResultForBlockIDRequest)
	if err != nil {
		return nil, err
	}

	flowExecResult, err := convert.MessageToExecutionResult(executionResultForBlockIDResponse.ExecutionResult)
	if err != nil {
		return nil, err
	}

	flowBlock, err := convert.MessageToBlock(blk)
	if err != nil {
		return nil, err
	}

	flowBlockStatus, err := convert.MessagesToBlockStatus(blockStatus)
	if err != nil {
		return nil, err
	}

	if err != nil {
		// handle case where execution result is not yet available
		if se, ok := status.FromError(err); ok {
			if se.Code() == codes.NotFound {
				err := block.Build(flowBlock, nil, link, flowBlockStatus, expandFields)
				if err != nil {
					return nil, err
				}
				return &block, nil
			}
		}
		return nil, err
	}

	err = block.Build(flowBlock, flowExecResult, link, flowBlockStatus, expandFields)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

type blockRequest struct {
	id     *flow.Identifier
	height uint64
	latest bool
	sealed bool
}

type blockRequestOption func(blkRequest *blockRequest)

func forID(id *flow.Identifier) blockRequestOption {
	return func(blockRequest *blockRequest) {
		blockRequest.id = id
	}
}
func forHeight(height uint64) blockRequestOption {
	return func(blockRequest *blockRequest) {
		blockRequest.height = height
	}
}

func forFinalized(queryParam uint64) blockRequestOption {
	return func(blockRequest *blockRequest) {
		switch queryParam {
		case request.SealedHeight:
			blockRequest.sealed = true
			fallthrough
		case request.FinalHeight:
			blockRequest.latest = true
		}
	}
}

// blockProvider is a layer of abstraction on top of the backend access.API and provides a uniform way to
// look up a block or a block header either by ID or by height
type blockRequestProvider struct {
	blockRequest
	backend access.API
}

func NewBlockRequestProvider(backend access.API, options ...blockRequestOption) *blockRequestProvider {
	blockRequestProvider := &blockRequestProvider{
		backend: backend,
	}

	for _, o := range options {
		o(&blockRequestProvider.blockRequest)
	}
	return blockRequestProvider
}

func (blkProvider *blockRequestProvider) getBlock(ctx context.Context) (*flow.Block, flow.BlockStatus, error) {
	if blkProvider.id != nil {
		blk, _, err := blkProvider.backend.GetBlockByID(ctx, *blkProvider.id)
		if err != nil { // unfortunately backend returns internal error status if not found
			return nil, flow.BlockStatusUnknown, NewNotFoundError(
				fmt.Sprintf("error looking up block with ID %s", blkProvider.id.String()), err,
			)
		}
		return blk, flow.BlockStatusUnknown, nil
	}

	if blkProvider.latest {
		blk, status, err := blkProvider.backend.GetLatestBlock(ctx, blkProvider.sealed)
		if err != nil {
			// cannot be a 'not found' error since final and sealed block should always be found
			return nil, flow.BlockStatusUnknown, NewRestError(http.StatusInternalServerError, "block lookup failed", err)
		}
		return blk, status, nil
	}

	blk, status, err := blkProvider.backend.GetBlockByHeight(ctx, blkProvider.height)
	if err != nil { // unfortunately backend returns internal error status if not found
		return nil, flow.BlockStatusUnknown, NewNotFoundError(
			fmt.Sprintf("error looking up block at height %d", blkProvider.height), err,
		)
	}
	return blk, status, nil
}

// blockProvider is a layer of abstraction on top of the accessproto.AccessAPIClient and provides a uniform way to
// look up a block or a block header either by ID or by height
type blockForwarderProvider struct {
	blockRequest
	upstream accessproto.AccessAPIClient
}

func NewBlockForwarderProvider(upstream accessproto.AccessAPIClient, options ...blockRequestOption) *blockForwarderProvider {
	blockForwarderProvider := &blockForwarderProvider{
		upstream: upstream,
	}

	for _, o := range options {
		o(&blockForwarderProvider.blockRequest)
	}
	return blockForwarderProvider
}

func (blkProvider *blockForwarderProvider) getBlock(ctx context.Context) (*entities.Block, entities.BlockStatus, error) {
	if blkProvider.id != nil {
		getBlockByIdRequest := &accessproto.GetBlockByIDRequest{
			Id: []byte(blkProvider.id.String()),
		}
		blockResponse, err := blkProvider.upstream.GetBlockByID(ctx, getBlockByIdRequest)
		if err != nil { // unfortunately grpc returns internal error status if not found
			return nil, entities.BlockStatus_BLOCK_UNKNOWN, NewNotFoundError(
				fmt.Sprintf("error looking up block with ID %s", blkProvider.id.String()), err,
			)
		}
		return blockResponse.Block, entities.BlockStatus_BLOCK_UNKNOWN, nil
	}

	if blkProvider.latest {
		getLatestBlockRequest := &accessproto.GetLatestBlockRequest{
			IsSealed: blkProvider.sealed,
		}
		blockResponse, err := blkProvider.upstream.GetLatestBlock(ctx, getLatestBlockRequest)
		if err != nil {
			// cannot be a 'not found' error since final and sealed block should always be found
			return nil, entities.BlockStatus_BLOCK_UNKNOWN, NewRestError(http.StatusInternalServerError, "block lookup failed", err)
		}
		return blockResponse.Block, blockResponse.BlockStatus, nil
	}

	getBlockByHeight := &accessproto.GetBlockByHeightRequest{
		Height:            blkProvider.height,
		FullBlockResponse: true,
	}
	blockResponse, err := blkProvider.upstream.GetBlockByHeight(ctx, getBlockByHeight)
	if err != nil { // unfortunately grpc returns internal error status if not found
		return nil, entities.BlockStatus_BLOCK_UNKNOWN, NewNotFoundError(
			fmt.Sprintf("error looking up block at height %d", blkProvider.height), err,
		)
	}
	return blockResponse.Block, blockResponse.BlockStatus, nil
}
