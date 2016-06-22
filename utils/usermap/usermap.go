// usermap.go
// This tool is for creating initial usermap contract
package main

import (
	"../../alg"
	"../../cli"
	"../../cmn"
	"../../solidity"
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	env_filename string = "alterorg.json"
)

func main() {

	if err := cmn.LoadSysEnv(env_filename); err != nil {
		fmt.Printf("error occured when loading sysenv file\n%s\n", err.Error())
		return
	}

	term := func() {
		fmt.Println("Terminating...")
		cli.TermEth(cli.STTS_ETH_NOT_START)
		time.Sleep(5 * time.Second)
	}
	cli.StartEth()
	defer cli.TermEth(cli.STTS_ETH_NOT_START)
	// TODO:will change to not to use sleep.
	defer time.Sleep(5 * time.Second)
	solidity.Init_usermap()
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
	go func() {
		for cli.GetEthStatus() != cli.STTS_ETH_STARTED {
			time.Sleep(1 * time.Second)
		}
		tx, err := alg.NewUserMap()
		if err != nil {
			fmt.Printf("error occured when creating UserMap\n%s\n", err.Error())
			term()
			return
		}
		address := ""
		for address == "" {
			receipt, err := cli.CheckContractTransaction(tx)
			if err != nil {
				fmt.Printf("error occured when checking transaction\n%s\n", err.Error())
				term()
				return
			}
			fmt.Printf("address:%s\n", receipt.CA)
			address = receipt.CA
			time.Sleep(1 * time.Second)
		}

	}()
	for scan.Scan() {
		if scan.Text() == "exit" {
			term()
			return
		}
	}

}
