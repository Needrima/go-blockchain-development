package block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	TimeStamp     int64    `json:"timestamp"`
	Nonce         int      `json:"nonce"`
	PrevHash      [32]byte `json:"previous_hash"`
	Transanctions []string `json:"transanctions"`
}

func NewBlock(nonce int, prevHash [32]byte) *Block {
	return &Block{
		TimeStamp: time.Now().UnixNano(),
		Nonce:     nonce,
		PrevHash:  prevHash,
	}
}

func (b *Block) Print() {
	fmt.Printf("Time stamp:                 %d\n\n", b.TimeStamp)
	fmt.Printf("Nonce:                		%d\n\n", b.Nonce)
	fmt.Printf("Previous hash:              %x\n\n", b.PrevHash)
	fmt.Printf("Transanctions:              %s\n\n", b.Transanctions)
}

func (b *Block) Hash() [32]byte {
	m, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	return sha256.Sum256(m)
}

type BlockChain struct {
	TransanctionsPool []string
	Blocks            []*Block
}

func NewBlockChain() *BlockChain {
	b := &Block{} // empty block to use as initial hash
	bc := &BlockChain{}
	bc.AddNewBlockToChain(0, b.Hash())

	return bc
}

func (bc *BlockChain) AddNewBlockToChain(nonce int, prevHash [32]byte) {
	block := NewBlock(nonce, prevHash)
	bc.Blocks = append(bc.Blocks, block)
}

func (bc *BlockChain) GetPreviousBlock() *Block {
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	return lastBlock
}

func (bc *BlockChain) Print() {
	for i, block := range bc.Blocks {
		fmt.Printf("%s Block %d %s\n", strings.Repeat("=", 10), i, strings.Repeat("=", 10))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 50))
}
