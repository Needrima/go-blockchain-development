package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PublicKey         *ecdsa.PublicKey
	PrivateKey        *ecdsa.PrivateKey
	BlockChainAddress string
}

func NewWallet() *Wallet {
	// 1. Creating ECDSA private key (32 bytes) public key (64 bytes)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	w := &Wallet{}
	w.PrivateKey = privateKey
	w.PublicKey = &privateKey.PublicKey

	// To get a wallet address follow the steps below (google search: Technical background of version 1 bitcoin address)

	// 2. Perform SHA-256 hashing on the public key (32 bytes).
	h2 := sha256.New()
	h2.Write(w.PublicKey.X.Bytes())
	h2.Write(w.PublicKey.Y.Bytes())
	digest2 := h2.Sum(nil)

	// 3. Perform RIPEMD-160 hashing on the result of SHA-256 (20 bytes).
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)

	// 4. Add version byte in front of RIPEMD-160 hash (0x00 for Main Network).
	vd4 := make([]byte, 21) // 20 bytes from SHA-256 with 1 byte from 0x00(version byte for main network)
	vd4[0] = 0x00
	copy(vd4[1:], digest3)

	// 5. Perform SHA-256 hash on the extended RIPEMD-160 result.
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	// 6. Perform SHA-256 hash on the result of the previous SHA-256 hash.
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	// 7. Take the first 4 bytes of the second SHA-256 hash for checksum.
	chsum := digest6[:4]

	// 8. Add the 4 checksum bytes from 7 at the end of extended RIPEMD-160 hash from 4 (25 bytes).
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4)
	copy(dc8[21:], chsum)

	// 9. Convert the result from a byte string into base58.
	address := base58.Encode(dc8)
	w.BlockChainAddress = address

	return w
}

func (w *Wallet) GetPrivateKey() *ecdsa.PrivateKey {
	return w.PrivateKey
}

func (w *Wallet) GetPrivateKeyString() string {
	return fmt.Sprintf("%x", w.PrivateKey.D.Bytes())
}

func (w *Wallet) GetPublicKey() *ecdsa.PublicKey {
	return w.PublicKey
}

func (w *Wallet) GetPublicKeyString() string {
	return fmt.Sprintf("%x%x", w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes())
}

func (w *Wallet) GetBlockChainAddress() string {
	return w.BlockChainAddress
}

type Transanction struct { 
	SenderPrivateKey           *ecdsa.PrivateKey
	SenderPublicKey            *ecdsa.PublicKey
	SenderBlockChainAddress    string
	RecipientBlockChainAddress string
	Value                      float32
}

type Signature struct {
	R *big.Int
	S *big.Int
}

func NewTransanction(senderPrivKey *ecdsa.PrivateKey, senderPubKey *ecdsa.PublicKey, senderBlockChainAddress, recipientBlockChainAddress string, value float32) *Transanction {
	return &Transanction{senderPrivKey, senderPubKey, senderBlockChainAddress, recipientBlockChainAddress, value}
}

func (t *Transanction) GenerateSignature() *Signature {
	// 1. Marshal the sender block chain address, recipient block chain address and value
	m, _ := json.Marshal(struct {
		Sender   string  `json:"sender_blockchain_address"`
		Receiver string  `json:"recipient_blockchain_address"`
		Value    float32 `json:"value"`
	}{
		Sender:   t.SenderBlockChainAddress,
		Receiver: t.RecipientBlockChainAddress,
		Value:    t.Value,
	})

	// 2. Hash the marshalled value
	hash := sha256.Sum256(m)

	// 3. Sign it and return signed value
	r, s, _ := ecdsa.Sign(rand.Reader, t.SenderPrivateKey, hash[:])

	return &Signature{R: r, S: s}
}

// adding a stringer method to make it printable as a string
func (s *Signature) String() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}
