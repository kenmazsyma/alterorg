// svr/api.go
// This source file enables access the functions of Alterorg
//
// Copyright holder is set forth in alterorg.md

package main

import (
	"./api"
	"./cli"
	"./solidity"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

/*
type Alterorg struct {
}

type Args struct {
	Hop string `json:"hop"`
	Hod string `json:"aaa"`
}

type Result int

func (self *Alterorg) Create(args Args, result *Result) error {
	fmt.Print("!!!!!!!!\n")
	fmt.Print("hop:" + args.Hop + "  hod:" + args.Hod)
	fmt.Print("\n!!!!!!!!\n")
	*result = 1
	return nil
}

func (self *Alterorg) Create2(a int, result *Result) error {
	fmt.Print("!!!!!!!!\n")
	fmt.Printf("test:%d", a)
	fmt.Print("\n!!!!!!!!\n")
	*result = 1
	return nil
}
*/
type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HttpConn) Close() error                      { return nil }

func main() {
	err := cli.InitEth("http://localhost:8545")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	solidity.Init_assembly()
	org := api.NewAssembly()
	sv := rpc.NewServer()
	sv.Register(org)
	l, _ := net.Listen("tcp", ":1234")
	defer l.Close()
	fmt.Print("received")
	http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCodec := jsonrpc.NewServerCodec(&HttpConn{in: r.Body, out: w})
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(200)
		err := sv.ServeRequest(serverCodec)
		if err != nil {
			fmt.Print(err)
			return
		}
	}))
}
