package execution

import (
	"github.com/dapperlabs/bamboo-emulator/data"
	"github.com/dapperlabs/bamboo-emulator/runtime"
)

// Node simulates the behaviour of a Bamboo execution node.
type Node struct {
	state    *data.WorldState
	computer *Computer
}

// NewNode returns a new simulated execution node.
func NewNode(state *data.WorldState) *Node {
	runtime := runtime.NewMockRuntime()
	computer := NewComputer(runtime, state.GetTransaction, state.GetRegister)

	return &Node{
		state:    state,
		computer: computer,
	}
}

// ExecuteBlock executes a block and saves the results to the world state.
func (n *Node) ExecuteBlock(block *data.Block) error {
	registers, results, err := n.computer.ExecuteBlock(block)
	if err != nil {
		return err
	}

	n.state.CommitRegisters(registers)
	err = n.state.SealBlock(block.Hash())
	if err != nil {
		return err
	}

	for hash, succeeded := range results {
		if succeeded {
			err = n.state.UpdateTransactionStatus(hash, data.TxSealed)
			if err != nil {
				return err
			}
		} else {
			err = n.state.UpdateTransactionStatus(hash, data.TxReverted)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
