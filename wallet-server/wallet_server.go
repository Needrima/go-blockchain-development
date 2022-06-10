package main

import (
	"encoding/json"
	"flag"
	"go-block-chain-dev/wallet"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type WalletServer struct {
	Port    uint16
	Gateway string
}

var walletCache = map[string]*wallet.Wallet{}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (ws *WalletServer) GetPort() uint16 {
	return ws.Port
}

func (ws *WalletServer) GetGateway() string {
	return ws.Gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("templates/index.html"))
	if err := tpl.Execute(w, nil); err != nil {
		log.Fatal("executing trmplate:", err)
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)

	http.HandleFunc("/wallet", func(w http.ResponseWriter, r *http.Request) {
		wall, ok := walletCache["wallet"]
		if !ok {
			wall = wallet.NewWallet()
			walletCache["wallet"] = wall
		}

		if err := json.NewEncoder(w).Encode(struct {
			Pub            string `json:"publicKey"`
			Priv           string `json:"privateKey"`
			BlockchainAddr string `json:"blockchainAddr"`
		}{
			Pub:            wall.GetPublicKeyString(),
			Priv:           wall.GetPrivateKeyString(),
			BlockchainAddr: wall.GetBlockChainAddress(),
		}); err != nil {
			log.Println("error occured", err)
		}
	})
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(ws.GetPort())), nil))
}

func main() {
	port := flag.Uint("port", 4682, "default port for wallet server")
	gateway := flag.String("gateway", "http://127.0.0.1:8080", "default port for blockchain server")
	flag.Parse()

	walletServer := NewWalletServer(uint16(*port), *gateway)

	walletServer.Run()
}
