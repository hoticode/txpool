// TxPool contains all currently known transactions.
// Transactions enter the pool when they are received from the network or submitted locally.
// They exit the pool when they are included in the blockchain.
// The pool seperate processable transactions (which can be applied to the current state) and future transactions.
// Transactions move between those two states over time as they are received and processed.

package core

import (
	"errors"
	"github.com/txpool/common"
	"github.com/txpool/core/types"
	"github.com/txpool/common/log"
)

// define errors description for error.
var (
	ErrDuplicateTxs = errors.New("duplicate transaction")
	ErrTxPoolHasFull = errors.New("txpool has full now")
	ErrAccountTxsPoolHasFull = errors.New("pending of the account has full now")
)

// DefaultTxPoolConfig contains the default configurations for the transaction pool.
var DefaultTxPoolConfig = TxPoolConfig{
	GlobalSlots:  4096,
}

// structure for tx lookup.
type txLookup struct {
	all  map[common.Hash]*types.Transaction
}

// newTxLookup returns a new txLookup structure.
func newTxLookup() *txLookup {
	return &txLookup{
		all: make(map[common.Hash]*types.Transaction),
	}
}

// TxPoolConfig are the configuration parameters of the transaction pool.
type TxPoolConfig struct {
	GlobalSlots  uint64       // Maximum number of executable transaction slots for txpool
}

// sanitize checks the provided user configurations and changes anything that's  unreasonable or unworkable.
func (config *TxPoolConfig) sanitize() TxPoolConfig {
	conf := *config
	if conf.GlobalSlots < 1 {
		log.Warn("Sanitizing invalid txpool global slots.")
		conf.GlobalSlots = DefaultTxPoolConfig.GlobalSlots
	}
	return conf
}

type TxsPool interface {
	// AddTx add a transaction to the txpool.
	AddTx(tx *types.Transaction) error

	// DelTxs delete the transactions which in processing queue.
	// Once a block was commited, transaction contained in the block can be removed.
	DelTxs() error

	// GetTxs gets the transactons which in pending status.
	GetTxs() []*types.Transaction
}

// As we desigin, tree queues which is all, pending and processing in txpool.
// The all receive transactions form network or local and reserve them to be scheduled by producer.
type TxPool struct {
	config       TxPoolConfig
	all          *txLookup
}

// NewTxPool creates a new transaction pool to gather, sort and filter inbound transactions from the network and local.
func NewTxPool(config TxPoolConfig) *TxPool {
	config = (&config).sanitize()

	// Create the transaction pool with its initial settings
	pool := &TxPool{
		config:      config,
		all:         newTxLookup(),
	}

	return pool
}

// Get returns a transaction if it exists in the lookup, or nil if not found.
func (t *txLookup) Get(hash common.Hash) *types.Transaction {
	return t.all[hash]
}

// Count returns the current number of items in the lookup.
func (t *txLookup) Count() int {
	return len(t.all)
}

// Add adds a transaction to the lookup.
func (t *txLookup) Add(tx *types.Transaction) {
	t.all[tx.Hash()] = tx
}

// Remove removes a transaction from the lookup.
func (t *txLookup) Remove(hash common.Hash) {
	delete(t.all, hash)
}

// Adding tx to the queue of all.
func (pool *TxPool) addToAll (tx *types.Transaction) {
	pool.all.Add(tx)
}

// Get pending txs from txpool.
func (pool *TxPool) GetTxs() []*types.Transaction {
	txList := make([]*types.Transaction, 0)
	for _, tx := range(pool.all.all){
		txList = append(txList, tx)
	}

	return txList
}

// Update processing queue, clean txs from process and all queue.
func (pool *TxPool) DelTxs(){
	// TODO: Adding a queue to sign some txs has been processed.
	log.Info("Update txpool after the txs has been applied by producer.")
	return
}

// Adding transaction to the txpool
func (pool *TxPool) AddTx(tx *types.Transaction) {
	if uint64(pool.all.Count()) > DefaultTxPoolConfig.GlobalSlots {
		log.Error("Txpool has full.")
		// TODO: return an sepcified error
		return
	}
	if nil != pool.all.Get(tx.Hash()){
		log.Error("The tx has exist, please confirm.")
		// TODO: return an sepcified error
		return
	}

	pool.all.Add(tx)
}