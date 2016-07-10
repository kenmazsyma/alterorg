// svr/svr.go
// This source file enables access the functions of Alterorg
//
// Copyright holder is set forth in alterorg.md

package main

import (
	"../alg"
	"../api"
	"../cli"
	"../cmn"
	"../solidity"
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HttpConn) Close() error                      { return nil }

const (
	env_filename string = "alterorg.json"
)

func main() {
	if err := cmn.LoadSysEnv(env_filename); err != nil {
		fmt.Printf("error occured when loading sysenv file\n%s\n", err.Error())
		return
	}
	if err := cmn.LoadApEnv(cmn.SysEnv.ApEnvPath); err != nil {
		fmt.Printf("error occured when loading apenv file\n%s\n", err.Error())
		return
	}
	cli.StartEth()
	cli.StartIpfs()
	defer cli.TermEth(cli.STTS_ETH_NOT_START)
	defer cli.TermIpfs(cli.STTS_IPFS_NOT_START)
	// TODO:will change to not to use sleep.
	defer time.Sleep(5 * time.Second)
	solidity.Init_assembly()
	solidity.Init_usermap()
	alg.UserMap_Prepare()
	sv := rpc.NewServer()
	sv.Register(api.NewAssembly())
	sv.Register(api.NewAlterorg())
	sv.Register(api.NewUser())
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Printf("Failed to run server:%s\n", err.Error())
	}
	defer l.Close()
	go http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCodec := jsonrpc.NewServerCodec(&HttpConn{in: r.Body, out: w})
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(200)
		err := sv.ServeRequest(serverCodec)
		if err != nil {
			fmt.Print(err)
			return
		}
	}))
	scan := bufio.NewScanner(os.Stdin)
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-signal_chan
		switch s {
		// kill -SIGHUP XXXX
		case syscall.SIGHUP:
			fmt.Println("hungup")

		// kill -SIGINT XXXX or Ctrl+c
		case syscall.SIGINT:
			fmt.Println("Warikomi")

		// kill -SIGTERM XXXX
		case syscall.SIGTERM:
			fmt.Println("force stop")
			return

		// kill -SIGQUIT XXXX
		case syscall.SIGQUIT:
			fmt.Println("stop and core dump")
			return

		default:
			fmt.Println("Unknown signal.")
			return
		}
	}()
	for scan.Scan() {
		if scan.Text() == "exit" {
			fmt.Println("Terminating...")
			return
		}
	}
}
