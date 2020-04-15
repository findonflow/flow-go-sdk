package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/examples"
	"github.com/dapperlabs/flow-go-sdk/keys"
	"github.com/dapperlabs/flow-go-sdk/templates"
)

func main() {
	AddAccountKeyDemo()
}

func AddAccountKeyDemo() {
	ctx := context.Background()
	accountAddr, accountKey, accountPrivateKey := examples.CreateAccount() // Creates a new account and returns the address+key

	newPrivateKey := examples.RandomPrivateKey() // Create the new key to add to your account
	newAccountKey := flow.AccountKey{
		PublicKey: newPrivateKey.PublicKey(),
		SignAlgo:  keys.ECDSA_P256_SHA3_256.SigningAlgorithm(),
		HashAlgo:  keys.ECDSA_P256_SHA3_256.HashingAlgorithm(),
		Weight:    keys.PublicKeyWeightThreshold,
	}

	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)

	// Create a Cadence script that will add another key to our account.
	addKeyScript, err := templates.AddAccountKey(newAccountKey)
	examples.Handle(err)

	// Create a transaction to execute the script.
	// The transaction is signed by our account key so it has permission to add keys.
	addKeyTx := flow.NewTransaction().
		SetScript(addKeyScript).
		SetProposalKey(accountAddr, accountKey.ID, accountKey.SequenceNumber).
		SetPayer(accountAddr).
		// This defines which accounts are accessed by this transaction
		AddAuthorizer(accountAddr)

	// Sign the transaction with the new account.
	err = addKeyTx.SignEnvelope(
		accountAddr,
		accountKey.ID,
		accountPrivateKey.Signer(),
	)
	examples.Handle(err)

	// Send the transaction to the network.
	err = flowClient.SendTransaction(ctx, *addKeyTx)
	examples.Handle(err)

	examples.WaitForSeal(ctx, flowClient, addKeyTx.ID())

	fmt.Println("Public key added to account!")

}