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

func Start() error {
	fmt.Println("Start")
	ethOutput := &appOutput{
		Callback: func(text string) {
			fmt.Println("eth:" + text)
		},
		write: func(text []byte) (int, error) {
			fmt.Printf("EthWrite:::%s\n", text)
			return 0, nil
		},
	}
	ipfsOutput := &appOutput{
		Callback: func(text string) {
			fmt.Println("ipfs:" + text)
		},
		write: func(text []byte) (int, error) {
			fmt.Printf("IpfsWrite:::%s\n", text)
			return 0, nil
		},
	}
	//prm := strings.Split(Env.EthPrm, " ")
	prm := splitArgs(Env.EthPrm)
	fmt.Printf("EthPrm:%s", strings.Join(prm, ":::"))
	ethCmd = exec.Command(Env.EthCmd, prm...)
	ethCmd.Stdout = ethOutput
	ethCmd.Stdin = &ethInput
	er := ethCmd.Start()
	if er != nil {
		fmt.Println(er.Error())
		return er
	}
	//prm = strings.Split(Env.IpfsPrm, " ")
	prm = splitArgs(Env.IpfsPrm)
	ipfsCmd = exec.Command(Env.IpfsCmd, prm...)
	ipfsCmd.Stdout = ipfsOutput
	ipfsCmd.Stdin = &ipfsInput
	er = ipfsCmd.Start()
	if er != nil {
		fmt.Println(er.Error())
		return er
	}
	// TODO:move url to env file & not use sleep
	time.Sleep(10 * time.Second)
	cli.InitEth(Env.EthUrl)
	cli.InitIpfs(Env.IpfsUrl)
	return nil
}

func Stop() error {
	//ethInput.WriteString("exit\n")
	//ipfsInput.WriteString("exit\n")
	//ethCmd.Process.Kill()
	ethCmd.Process.Signal(syscall.SIGINT)
	ipfsCmd.Process.Signal(syscall.SIGINT)
	//ipfsCmd.Process.Kill()
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
