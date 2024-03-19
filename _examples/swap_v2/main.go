package main

import (
	"context"
	"fmt"
	"github.com/Rodert/jupiter-go/jupiter"
	"github.com/Rodert/jupiter-go/solana"
	"time"
)

const (
	UserPublicKey    = "your public key here"
	walletPrivateKey = "your private key"
)

func main() {
	swap := GetSwapJson()

	ctx := context.TODO()
	// ... previous code
	// swap := swapResponse.JSON200

	// Create a wallet from private key
	wallet, err := solana.NewWalletFromPrivateKeyBase58(walletPrivateKey)
	if err != nil {
		// handle me
	}

	// Create a Solana client
	solanaClient, err := solana.NewClient(wallet, "https://api.mainnet-beta.solana.com")
	if err != nil {
		// handle me
	}

	// Sign and send the transaction
	signedTx, err := solanaClient.SendTransactionOnChain(ctx, swap.SwapTransaction)
	if err != nil {
		// handle me
	}

	// wait a bit to let the transaction propagate to the network
	// this is just an example and not a best practice
	// you could use a ticker or wait until we implement the WebSocket monitoring ;)
	time.Sleep(20 * time.Second)

	// Get the status of the transaction (pull the status from the blockchain at intervals
	// until the transaction is confirmed)
	confirmed, err := solanaClient.CheckSignature(ctx, signedTx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("confirmed: %+v\n", confirmed)
}

func GetSwapJson() *jupiter.SwapResponse {
	jupClient, err := jupiter.NewClientWithResponses(jupiter.DefaultAPIURL)
	if err != nil {
		panic(err)
	}

	ctx := context.TODO()

	slippageBps := 250

	// Get the current quote for a swap
	quoteResponse, err := jupClient.GetQuoteWithResponse(ctx, &jupiter.GetQuoteParams{
		InputMint:   "So11111111111111111111111111111111111111112",
		OutputMint:  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		Amount:      10000,
		SlippageBps: &slippageBps,
	})
	if err != nil {
		panic(err)
	}

	if quoteResponse.JSON200 == nil {
		panic("invalid GetQuoteWithResponse response")
	}

	quote := quoteResponse.JSON200

	// More info: https://station.jup.ag/docs/apis/troubleshooting
	prioritizationFeeLamports := jupiter.SwapRequest_PrioritizationFeeLamports{}
	if err = prioritizationFeeLamports.UnmarshalJSON([]byte(`"auto"`)); err != nil {
		panic(err)
	}

	dynamicComputeUnitLimit := true
	// Get instructions for a swap
	swapResponse, err := jupClient.PostSwapWithResponse(ctx, jupiter.PostSwapJSONRequestBody{
		PrioritizationFeeLamports: &prioritizationFeeLamports,
		QuoteResponse:             *quote,
		UserPublicKey:             UserPublicKey,
		DynamicComputeUnitLimit:   &dynamicComputeUnitLimit,
	})
	if err != nil {
		panic(err)
	}

	if swapResponse.JSON200 == nil {
		panic("invalid PostSwapWithResponse response")
	}

	swap := swapResponse.JSON200
	fmt.Printf("%+v", swap)
	return swap
}
