package main

import (
	"fmt"
	"go-block-chain-dev/wallet"
)

func main() {
	w := wallet.NewWallet()

	transanction := wallet.NewTransanction(w.GetPrivateKey(), w.GetPublicKey(), w.BlockChainAddress, "recipient blockcahin address", 2.0)
	fmt.Println("signature:", transanction.GenerateSignature().String())
}
