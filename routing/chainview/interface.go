package chainview

import (
	"github.com/roasbeef/btcd/chaincfg/chainhash"
	"github.com/roasbeef/btcd/wire"
)

// FilteredChainView repersents a subscription to a certain subset of of the
// UTXO set for a particular chain. This interface is useful from the point of
// view of maintaining an up-to-date channel graph for the Lighting Network.
// The subset of the UTXO to be subscribed is that of all the currently opened
// channels. Each time a channel is closed (the output is spent), a
// notification is to be sent allowing the graph to be pruned.
//
// NOTE: As FilteredBlocks are generated, it is recommended that
// implementations reclaim the space occupied by newly spent outputs.
type FilteredChainView interface {
	// FilteredBlocks returns the channel that filtered blocks are to be
	// sent over. Each time a block is connected to the end of a main
	// chain, and appropriate FilteredBlock which contains the transactions
	// which mutate our watched UTXO set is to be returned.
	FilteredBlocks() <-chan *FilteredBlock

	// DisconnectedBlocks returns a receive only channel which will be sent
	// upon with the empty filtered blocks of blocks which are disconnected
	// from the main chain in the case of a re-org.
	DisconnectedBlocks() <-chan *FilteredBlock

	// UpdateFilter updates the UTXO filter which is to be consulted when
	// creating FilteredBlocks to be sent to subscribed clients. This
	// method is cumulative meaning repeated calls to this method should
	// _expand_ the size of the UTXO sub-set currently being watched.  If
	// the set updateHeight is _lower_ than the best known height of the
	// implementation, then the state should be rewound to ensure all
	// relevant notifications are dispatched.
	UpdateFilter(ops []wire.OutPoint, updateHeight uint32) error

	// FilterBlock takes a block hash, and returns a FilteredBlocks which
	// is the result of applying the current registered UTXO sub-set on the
	// block corresponding to that block hash.
	//
	// TODO(roasbeef): make a version that does by height also?
	FilterBlock(blockHash *chainhash.Hash) (*FilteredBlock, error)

	// Start starts all goroutine necessary for the operation of the
	// FilteredChainView implementation.
	Start() error

	// Stop stops all goroutines which we launched by the prior call to the
	// Start method.
	Stop() error
}

// FilteredBlock is a block which includes the transactions that modify the
// subscribed sub-set of the UTXO set registered to the current
// FilteredChainView concrete implementation.
type FilteredBlock struct {
	// Hash is the hash of the newly filtered block.
	Hash chainhash.Hash

	// Height is the height of the newly filtered block.
	Height uint32

	// Transactions is the set of transactions which modify (spend) the
	// subscribed UTXO subset.
	Transactions []*wire.MsgTx
}
