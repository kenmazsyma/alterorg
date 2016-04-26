package cmn

import (
	//	"encoding/json"
	///	"errors"
	"fmt"
	//	"io/ioutil"
	//	"os"
	"../cli"
	"bytes"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var ethInput bytes.Buffer
var ipfsInput bytes.Buffer
var ethCmd *exec.Cmd
var ipfsCmd *exec.Cmd

type appOutput struct {
	Callback func(string)
	write    func([]byte) (int, error)
}

func (self *appOutput) Write(text []byte) (int, error) {
	return self.write(text)
}

func Initialize() error {
	return nil
}

const (
	NOUSE = iota
	RUN
	WAIT
	ERROR
)

var EthState int
var IpfsState int

func Start() error {
	if SysEnv.EthRun != 0 {
		ethOutput := &appOutput{
			Callback: func(text string) {
				fmt.Println("eth:" + text)
			},
			write: func(text []byte) (int, error) {
				fmt.Printf("EthWrite:::%s\n", text)
				return 0, nil
			},
		}
		prm := splitArgs(SysEnv.EthPrm)
		fmt.Printf("EthPrm:%s", strings.Join(prm, ":::"))
		ethCmd = exec.Command(SysEnv.EthCmd, prm...)
		ethCmd.Stdout = ethOutput
		ethCmd.Stdin = &ethInput
		er := ethCmd.Start()
		if er != nil {
			fmt.Println(er.Error())
			return er
		}
	}
	if SysEnv.IpfsRun != 0 {
		ipfsOutput := &appOutput{
			Callback: func(text string) {
				fmt.Println("ipfs:" + text)
			},
			write: func(text []byte) (int, error) {
				fmt.Printf("IpfsWrite:::%s\n", text)
				return 0, nil
			},
		}
		prm := splitArgs(SysEnv.IpfsPrm)
		ipfsCmd = exec.Command(SysEnv.IpfsCmd, prm...)
		ipfsCmd.Stdout = ipfsOutput
		ipfsCmd.Stdin = &ipfsInput
		er := ipfsCmd.Start()
		if er != nil {
			fmt.Println(er.Error())
			return er
		}
	}
	// TODO:move url to env file & not use sleep
	//time.Sleep(10 * time.Second)
	if SysEnv.EthRun != 0 {
		cli.InitEth(SysEnv.EthUrl)
		EthState = WAIT
		go func() {
			for i := 0; i < 10; i++ {
				time.Sleep(1 * time.Second)
				start, er := cli.GetEthListening()
				if er != nil {
					EthState = ERROR
				} else {
					if start == true {
						EthState = RUN
						fmt.Printf("Eth:Run:\n")
						break
					}
				}
			}
			if EthState == RUN {
				cli.GetCoinbase()
			}
		}()
	} else {
		EthState = NOUSE
	}
	if SysEnv.IpfsRun != 0 {
		IpfsState = RUN
		cli.InitIpfs(SysEnv.IpfsUrl)
	} else {
		IpfsState = NOUSE
	}
	return nil
}

func Stop() error {
	//if ethCmd != nil && ethCmd.Process != nil {
	if EthState == RUN {
		ethCmd.Process.Signal(syscall.SIGINT)
	}
	//if ipfsCmd != nil && ipfsCmd.Process != nil {
	if IpfsState == RUN {
		ipfsCmd.Process.Signal(syscall.SIGINT)
	}
	return nil
}

func splitArgs(txt string) []string {
	prm := strings.Split(txt, " ")
	for i := range prm {
		start := 0
		end := len(prm[i])
		if prm[i][0] == '"' {
			start = 1
			end--
		}
		//	if prm[i][end-1] == '"' {
		//		end--
		//	}
		prm[i] = prm[i][start:end]
	}
	return prm
}
