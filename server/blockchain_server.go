package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type BlockchainServer struct {
	Port uint16
}

func New(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bcs *BlockchainServer) GetPort() uint16 {
	return bcs.Port
}

func (bcs *BlockchainServer) Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}

func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/", bcs.Hello)

	http.ListenAndServe(":" + strconv.Itoa(int(bcs.GetPort())), nil)
}

func main() {
	s := New(8080)
	s.Run()
}
