// api.go
// This source file enables access the functions of Alterorg
//
// Copyright holder is set forth in alterorg.md

package main

import (
	"./api"
	"os/signal"
	"syscall"
	"time"
	//	"./cli"
	"./cmn"
	"./solidity"
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HttpConn) Close() error                      { return nil }

func main() {
	err := cmn.LoadSysEnv("alterorg.json")
	if err != nil {
		fmt.Printf("error occured when loading sysenv file\n%s\n", err.Error())
		return
	}
	// TODO:
	err = cmn.LoadApEnv(cmn.SysEnv.ApEnvPath)
	if err != nil {
		fmt.Printf("error occured when loading apenv file\n%s\n", err.Error())
		return
	}
	err = cmn.Start()
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	solidity.Init_assembly()
	sv := rpc.NewServer()
	sv.Register(api.NewAssembly())
	sv.Register(api.NewAlterorg())
	l, _ := net.Listen("tcp", ":1234")
	defer l.Close()
	fmt.Print("received")
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
			cmn.Stop()
			// TODO:don't use sleep.
			time.Sleep(5 * time.Second)
			return
		}
	}
}
