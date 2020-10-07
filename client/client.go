package client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/MariusVanDerWijden/ShareMyRPC/raiden"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	generalCost         = 4
	stateAccessCost     = 4
	filterCost          = 8
	contractCallingCost = 3
	estimateGasCost     = 1
	sendTransactionCost = 1
)

// Client is a thin wrapper around an ethclient object.
type Client struct {
	c     *ethclient.Client
	r     *raiden.Raiden
	token string
	other string
}

func NewClient(c *ethclient.Client, r *raiden.Raiden) *Client {
	return &Client{
		c: c,
		r: r,
	}
}

func NewClientFromURL(clientURL, raidenURL string) (*Client, error) {
	c, err := ethclient.Dial(clientURL)
	if err != nil {
		return nil, err
	}
	r := raiden.NewRaiden(raidenURL)
	return NewClient(c, r), nil
}

// Close closes the connection to the peer.
func (ec *Client) Close() {
	ec.c.Close()
}

// Send sends some money to the peer.
func (ec *Client) Send(amount int) {
	am := fmt.Sprint("%v", amount)
	if err := ec.r.PayToken(ec.token, ec.other, am); err != nil {
		panic(err)
	}
}

// ChainID retrieves the current chain ID for transaction replay protection.
func (ec *Client) ChainID(ctx context.Context) (*big.Int, error) {
	return ec.c.ChainID(ctx)
}

// BlockByHash returns the given full block.
func (ec *Client) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	ec.Send(generalCost)
	return ec.c.BlockByHash(ctx, hash)
}

// BlockByNumber returns a block from the current canonical chain. If number is nil, the
// latest known block is returned.
func (ec *Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	ec.Send(generalCost)
	return ec.c.BlockByNumber(ctx, number)
}

// BlockNumber returns the most recent block number
func (ec *Client) BlockNumber(ctx context.Context) (uint64, error) {
	ec.Send(generalCost)
	return ec.c.BlockNumber(ctx)
}

// HeaderByHash returns the block header with the given hash.
func (ec *Client) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	ec.Send(generalCost)
	return ec.c.HeaderByHash(ctx, hash)
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (ec *Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	ec.Send(generalCost)
	return ec.c.HeaderByNumber(ctx, number)
}

// TransactionByHash returns the transaction with the given hash.
func (ec *Client) TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	ec.Send(generalCost)
	return ec.c.TransactionByHash(ctx, hash)
}

// TransactionSender returns the sender address of the given transaction.
func (ec *Client) TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error) {
	ec.Send(generalCost)
	return ec.c.TransactionSender(ctx, tx, block, index)
}

// TransactionCount returns the total number of transactions in the given block.
func (ec *Client) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	ec.Send(generalCost)
	return ec.c.TransactionCount(ctx, blockHash)
}

// TransactionInBlock returns a single transaction at index in the given block.
func (ec *Client) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	ec.Send(generalCost)
	return ec.c.TransactionInBlock(ctx, blockHash, index)
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (ec *Client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	ec.Send(generalCost)
	return ec.c.TransactionReceipt(ctx, txHash)
}

// SyncProgress retrieves the current progress of the sync algorithm. If there's
// no sync currently running, it returns nil.
func (ec *Client) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	ec.Send(generalCost)
	return ec.c.SyncProgress(ctx)
}

// SubscribeNewHead subscribes to notifications about the current blockchain head
// on the given channel.
func (ec *Client) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	ec.Send(generalCost)
	return ec.c.SubscribeNewHead(ctx, ch)
}

// State Access

// NetworkID returns the network ID (also known as the chain ID) for this chain.
func (ec *Client) NetworkID(ctx context.Context) (*big.Int, error) {
	ec.Send(stateAccessCost)
	return ec.c.NetworkID(ctx)
}

// BalanceAt returns the wei balance of the given account.
// The block number can be nil, in which case the balance is taken from the latest known block.
func (ec *Client) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	ec.Send(stateAccessCost)
	return ec.c.BalanceAt(ctx, account, blockNumber)
}

// StorageAt returns the value of key in the contract storage of the given account.
// The block number can be nil, in which case the value is taken from the latest known block.
func (ec *Client) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	ec.Send(stateAccessCost)
	return ec.c.StorageAt(ctx, account, key, blockNumber)
}

// CodeAt returns the contract code of the given account.
// The block number can be nil, in which case the code is taken from the latest known block.
func (ec *Client) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	ec.Send(stateAccessCost)
	return ec.c.CodeAt(ctx, account, blockNumber)
}

// NonceAt returns the account nonce of the given account.
// The block number can be nil, in which case the nonce is taken from the latest known block.
func (ec *Client) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	ec.Send(stateAccessCost)
	return ec.c.NonceAt(ctx, account, blockNumber)
}

// Filters

// FilterLogs executes a filter query.
func (ec *Client) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	ec.Send(filterCost)
	return ec.c.FilterLogs(ctx, q)
}

// SubscribeFilterLogs subscribes to the results of a streaming filter query.
func (ec *Client) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	ec.Send(filterCost)
	return ec.c.SubscribeFilterLogs(ctx, q, ch)
}

// Pending State

// PendingBalanceAt returns the wei balance of the given account in the pending state.
func (ec *Client) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	ec.Send(stateAccessCost)
	return ec.c.PendingBalanceAt(ctx, account)
}

// PendingStorageAt returns the value of key in the contract storage of the given account in the pending state.
func (ec *Client) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error) {
	ec.Send(stateAccessCost)
	return ec.c.PendingStorageAt(ctx, account, key)
}

// PendingCodeAt returns the contract code of the given account in the pending state.
func (ec *Client) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	ec.Send(stateAccessCost)
	return ec.c.PendingCodeAt(ctx, account)
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (ec *Client) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	ec.Send(stateAccessCost)
	return ec.c.PendingNonceAt(ctx, account)
}

// PendingTransactionCount returns the total number of transactions in the pending state.
func (ec *Client) PendingTransactionCount(ctx context.Context) (uint, error) {
	ec.Send(stateAccessCost)
	return ec.c.PendingTransactionCount(ctx)
}

// Contract Calling

// CallContract executes a message call transaction, which is directly executed in the VM
// of the node, but never mined into the blockchain.
func (ec *Client) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	ec.Send(contractCallingCost)
	return ec.c.CallContract(ctx, msg, blockNumber)
}

// PendingCallContract executes a message call transaction using the EVM.
// The state seen by the contract call is the pending state.
func (ec *Client) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	ec.Send(contractCallingCost)
	return ec.c.PendingCallContract(ctx, msg)
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (ec *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	ec.Send(estimateGasCost)
	return ec.c.SuggestGasPrice(ctx)
}

// EstimateGas tries to estimate the gas needed to execute a specific transaction based on
// the current pending state of the backend blockchain. There is no guarantee that this is
// the true gas limit requirement as other transactions may be added or removed by miners,
// but it should provide a basis for setting a reasonable default.
func (ec *Client) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	ec.Send(estimateGasCost)
	return ec.c.EstimateGas(ctx, msg)
}

// SendTransaction injects a signed transaction into the pending pool for execution.
//
// If the transaction was a contract creation use the TransactionReceipt method to get the
// contract address after the transaction has been mined.
func (ec *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	ec.Send(sendTransactionCost)
	return ec.c.SendTransaction(ctx, tx)
}
