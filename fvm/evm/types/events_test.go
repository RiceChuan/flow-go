package types_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/onflow/cadence"

	"github.com/onflow/flow-go/fvm/systemcontracts"
	"github.com/onflow/flow-go/model/flow"

	"github.com/onflow/go-ethereum/core/vm"

	"github.com/onflow/cadence/encoding/ccf"
	cdcCommon "github.com/onflow/cadence/runtime/common"
	gethCommon "github.com/onflow/go-ethereum/common"
	gethTypes "github.com/onflow/go-ethereum/core/types"
	"github.com/onflow/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/fvm/evm/testutils"
	"github.com/onflow/flow-go/fvm/evm/types"
)

var evmLocation = cdcCommon.NewAddressLocation(
	nil,
	cdcCommon.Address(systemcontracts.SystemContractsForChain(flow.Emulator).EVMContract.Address),
	"EVM",
)

func TestEVMBlockExecutedEventCCFEncodingDecoding(t *testing.T) {
	t.Parallel()

	block := &types.Block{
		Height:          2,
		Timestamp:       100,
		TotalSupply:     big.NewInt(1500),
		ParentBlockHash: gethCommon.HexToHash("0x2813452cff514c3054ac9f40cd7ce1b016cc78ab7f99f1c6d49708837f6e06d1"),
		ReceiptRoot:     gethCommon.Hash{},
		TotalGasUsed:    15,
		TransactionHashes: []gethCommon.Hash{
			gethCommon.HexToHash("0x70b67ce6710355acf8d69b2ea013d34e212bc4824926c5d26f189c1ca9667246"),
		},
	}

	event := types.NewBlockEvent(block)
	ev, err := event.Payload.ToCadence(evmLocation)
	require.NoError(t, err)

	bep, err := types.DecodeBlockEventPayload(ev)
	require.NoError(t, err)

	assert.Equal(t, bep.Height, block.Height)

	blockHash, err := block.Hash()
	require.NoError(t, err)
	assert.Equal(t, bep.Hash, types.BytesToCadenceUInt8ArrayValue(blockHash.Bytes()))

	assert.Equal(t, bep.TotalSupply.Value, block.TotalSupply)
	assert.Equal(t, bep.Timestamp, block.Timestamp)
	assert.Equal(t, bep.TotalGasUsed, block.TotalGasUsed)
	assert.Equal(t, bep.ParentBlockHash, types.BytesToCadenceUInt8ArrayValue(block.ParentBlockHash.Bytes()))
	assert.Equal(t, bep.ReceiptRoot, types.BytesToCadenceUInt8ArrayValue(block.ReceiptRoot.Bytes()))

	hashes := make([]cadence.Array, len(block.TransactionHashes))
	for i, h := range block.TransactionHashes {
		hashes[i] = types.BytesToCadenceUInt8ArrayValue(h.Bytes())
	}
	assert.Equal(t,
		bep.TransactionHashes,
		hashes,
	)

	v, err := ccf.Encode(ev)
	require.NoError(t, err)
	assert.Equal(t, ccf.HasMsgPrefix(v), true)

	evt, err := ccf.Decode(nil, v)
	require.NoError(t, err)

	sc := systemcontracts.SystemContractsForChain(flow.Emulator)

	assert.Equal(t,
		cdcCommon.NewAddressLocation(
			nil,
			cdcCommon.Address(sc.EVMContract.Address),
			string(types.EventTypeBlockExecuted),
		).ID(),
		evt.Type().ID(),
	)
}

func TestEVMTransactionExecutedEventCCFEncodingDecoding(t *testing.T) {
	t.Parallel()

	txEncoded := "fff83b81ff0194000000000000000000000000000000000000000094000000000000000000000000000000000000000180895150ae84a8cdf00000825208"
	txBytes, err := hex.DecodeString(txEncoded)
	require.NoError(t, err)
	txHash := testutils.RandomCommonHash(t)
	blockHash := testutils.RandomCommonHash(t)
	data := "000000000000000000000000000000000000000000000000000000000000002a"
	dataBytes, err := hex.DecodeString(data)
	require.NoError(t, err)
	blockHeight := uint64(2)
	deployedAddress := types.NewAddress(gethCommon.HexToAddress("0x99466ed2e37b892a2ee3e9cd55a98b68f5735db2"))
	log := &gethTypes.Log{
		Index:       1,
		BlockNumber: blockHeight,
		BlockHash:   blockHash,
		TxHash:      txHash,
		TxIndex:     3,
		Address:     gethCommon.HexToAddress("0x99466ed2e37b892a2ee3e9cd55a98b68f5735db2"),
		Data:        dataBytes,
		Topics: []gethCommon.Hash{
			gethCommon.HexToHash("0x24abdb5865df5079dcc5ac590ff6f01d5c16edbc5fab4e195d9febd1114503da"),
		},
	}
	vmError := vm.ErrOutOfGas
	txResult := &types.Result{
		VMError:                 vmError,
		TxType:                  255,
		GasConsumed:             23200,
		DeployedContractAddress: &deployedAddress,
		ReturnedData:            dataBytes,
		Logs:                    []*gethTypes.Log{log},
		TxHash:                  txHash,
	}

	t.Run("evm.TransactionExecuted with failed status", func(t *testing.T) {
		event := types.NewTransactionEvent(txResult, txBytes, blockHeight)
		ev, err := event.Payload.ToCadence(evmLocation)
		require.NoError(t, err)

		tep, err := types.DecodeTransactionEventPayload(ev)
		require.NoError(t, err)

		assert.Equal(t, tep.BlockHeight, blockHeight)
		assert.Equal(t, tep.Hash, types.BytesToCadenceUInt8ArrayValue(txHash.Bytes()))
		assert.Equal(t, tep.Payload, types.BytesToCadenceUInt8ArrayValue(txBytes))
		assert.Equal(t, types.ErrorCode(tep.ErrorCode), types.ExecutionErrCodeOutOfGas)
		assert.Equal(t, tep.TransactionType, txResult.TxType)
		assert.Equal(t, tep.GasConsumed, txResult.GasConsumed)
		assert.Equal(t, tep.ErrorMessage, txResult.VMError.Error())
		assert.Equal(t, tep.ReturnedData, types.BytesToCadenceUInt8ArrayValue(txResult.ReturnedData))
		assert.Equal(
			t,
			tep.ContractAddress,
			txResult.DeployedContractAddress.ToCommon().Hex(),
		)

		encodedLogs, err := rlp.EncodeToBytes(txResult.Logs)
		require.NoError(t, err)
		assert.Equal(t, tep.Logs, types.BytesToCadenceUInt8ArrayValue(encodedLogs))

		v, err := ccf.Encode(ev)
		require.NoError(t, err)
		assert.Equal(t, ccf.HasMsgPrefix(v), true)

		evt, err := ccf.Decode(nil, v)
		require.NoError(t, err)

		assert.Equal(t, evt.Type().ID(), fmt.Sprintf(
			"A.%s.EVM.TransactionExecuted",
			systemcontracts.SystemContractsForChain(flow.Emulator).EVMContract.Address,
		))
	})

	t.Run("evm.TransactionExecuted with non-failed status", func(t *testing.T) {
		txResult.VMError = nil

		event := types.NewTransactionEvent(txResult, txBytes, blockHeight)
		ev, err := event.Payload.ToCadence(evmLocation)
		require.NoError(t, err)

		tep, err := types.DecodeTransactionEventPayload(ev)
		require.NoError(t, err)

		assert.Equal(t, tep.BlockHeight, blockHeight)
		assert.Equal(t, tep.Hash, types.BytesToCadenceUInt8ArrayValue(txHash.Bytes()))
		assert.Equal(t, tep.Payload, types.BytesToCadenceUInt8ArrayValue(txBytes))
		assert.Equal(t, types.ErrCodeNoError, types.ErrorCode(tep.ErrorCode))
		assert.Equal(t, tep.TransactionType, txResult.TxType)
		assert.Equal(t, tep.GasConsumed, txResult.GasConsumed)
		assert.Empty(t, tep.ErrorMessage)
		assert.Equal(t, tep.ReturnedData, types.BytesToCadenceUInt8ArrayValue(txResult.ReturnedData))
		assert.NotNil(t, txResult.DeployedContractAddress)
		assert.Equal(
			t,
			tep.ContractAddress,
			txResult.DeployedContractAddress.ToCommon().Hex(),
		)

		encodedLogs, err := rlp.EncodeToBytes(txResult.Logs)
		require.NoError(t, err)
		assert.Equal(t, tep.Logs, types.BytesToCadenceUInt8ArrayValue(encodedLogs))

		v, err := ccf.Encode(ev)
		require.NoError(t, err)
		assert.Equal(t, ccf.HasMsgPrefix(v), true)

		evt, err := ccf.Decode(nil, v)
		require.NoError(t, err)

		assert.Equal(t, evt.Type().ID(), fmt.Sprintf(
			"A.%s.EVM.TransactionExecuted",
			systemcontracts.SystemContractsForChain(flow.Emulator).EVMContract.Address,
		))
	})
}
