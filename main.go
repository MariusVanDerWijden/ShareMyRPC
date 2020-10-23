package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/MariusVanDerWijden/ShareMyRPC/client"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

var (
	raidenURL = "http://localhost:5001/api/v1"
	nodeURL   = "http://127.0.0.1:8546"
	token     = "0x95B2d84De40a0121061b105E6B54016a49621B44"
	peer      = "0x1F916ab5cf1B30B22f24Ebf435f53Ee665344Acf"
	wrongPeer = "0x0000000000000000000000000000000000000000"

	SK = "0xcdfbe6f7602f67a97602e3e9fc24cde1cdffa88acd47745c0b84c5ff55891e1b"
	sk = crypto.ToECDSAUnsafe(common.FromHex(SK))
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		runDemo(context.TODO())
	} else {
		demoFail(context.TODO())
	}
}

func runDemo(ctx context.Context) {
	fmt.Printf("Connecting to the nodes at geth-node: %v raiden-node: %v\n", nodeURL, raidenURL)
	client, err := client.NewClientFromURL(nodeURL, raidenURL)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Initializing node with peer: %v token %v\n", peer, token)
	client.Init(token, peer)
	// Run some basic commands
	queryBlockNum(ctx, client)
	estimateGas(ctx, client)
	sendTransaction(ctx, client)
}

func queryBlockNum(ctx context.Context, client *client.Client) {
	fmt.Println("Query current block number")
	no, err := client.BlockNumber(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Highest block is %v\n", no)
}

func estimateGas(ctx context.Context, client *client.Client) {
	from, to := common.HexToAddress("0x00012345"), common.HexToAddress("0x00054321")
	msg := ethereum.CallMsg{
		From: from,
		To:   &to,
		Data: []byte{1, 2, 3, 4, 5, 6},
	}
	fmt.Printf("Estimating gas for tx: %v\n", msg)
	gas, err := client.EstimateGas(ctx, msg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Estimated gas as %v\n", gas)
}

func sendTransaction(ctx context.Context, client *client.Client) {
	tx := createTx(ctx, client)
	fmt.Printf("Sending transaction: %v\n", tx)
	err := client.SendTransaction(ctx, tx)
	if err != nil {
		fmt.Println(err)
	}
	if err.Error() == "EOF" {
		fmt.Println("Sending transaction failed")
	} else {
		fmt.Println("Sending transaction successful")
	}
}

func createTx(ctx context.Context, client *client.Client) *types.Transaction {
	addr := crypto.PubkeyToAddress(sk.PublicKey)
	fmt.Printf("Retrieve the pending nonce for account: %v\n", addr)
	nonce, err := client.NonceAt(ctx, addr, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Pending nonce: %v\n", nonce)
	to := common.HexToAddress("0xABCD")
	amount := big.NewInt(10 * params.GWei)
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(10 * params.GWei)
	data := []byte{}
	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(nil), sk)
	if err != nil {
		fmt.Println(err)
	}
	return signedTx
}

func demoFail(ctx context.Context) {
	fmt.Printf("Connecting to the nodes at geth-node: %v raiden-node: %v\n", nodeURL, raidenURL)
	client, err := client.NewClientFromURL(nodeURL, raidenURL)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Initializing node with peer: %v token %v\n", wrongPeer, token)
	client.Init(token, wrongPeer)
	// Run some basic commands
	queryBlockNum(ctx, client)
	estimateGas(ctx, client)
	sendTransaction(ctx, client)
}
