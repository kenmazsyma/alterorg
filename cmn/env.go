package cmn

import (
	//	"encoding/json"
	///	"errors"
	"fmt"
	//	"io/ioutil"
	"../alg"
	"../cli"
	"bytes"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var ethInput bytes.Buffer
var ipfsInput bytes.Buffer
var ethCmd *exec.Cmd
var ipfsCmd *exec.Cmd
var Iwasmined bool

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
		/*ethOutput := &appOutput{
			Callback: func(text string) {
				fmt.Println("eth:" + text)
			},
			write: func(text []byte) (int, error) {
				fmt.Printf("EthWrite:::%s\n", text)
				return 0, nil
			},
		}*/
		/*ethOutput2 := &appOutput{
			Callback: func(text string) {
				fmt.Println("eth:" + text)
			},
			write: func(text []byte) (int, error) {
				fmt.Printf("EthError:::%s\n", text)
				return 0, nil
			},
		}*/
		prm := splitArgs(SysEnv.EthPrm)
		fmt.Printf("EthPrm:%s", strings.Join(prm, ":::"))
		ethCmd = exec.Command(SysEnv.EthCmd, prm...)
		ethCmd.Stdout = os.Stdout //ethOutput
		//ethCmd.Stderr = ethOutput2
		//	ethCmd.CombinedOutput = func(text []byte) (int, error) {
		//		fmt.Printf("EthWrite:::%s\n", text)
		//		return 0, nil
		//	}
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
	if SysEnv.EthRun != 0 {
		cli.InitEth(SysEnv.EthUrl)
		EthState = WAIT
		go func() {
			run := false
			for i := 0; i < 10; i++ {
				time.Sleep(1 * time.Second)
				start, er := cli.GetEthListening()
				if er != nil {
					EthState = ERROR
				} else {
					if start == true {
						// TODO:now searching a way of getting unlock timing
						// if delete below sleep, "account is locked" error occurs
						// because unlock proc in geth is not called yet at this timing.
						time.Sleep(2 * time.Second)
						run = true
						fmt.Printf("Eth:Run:\n")
						break
					}
				}
			}
			if run {
				time.Sleep(5 * time.Second)
				initEnv()
				//cli.GetCoinbase()
				//cli.GetUsrSet()
				EthState = RUN
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

func IsStart() bool {
	if (SysEnv.EthRun == 0 || EthState == RUN) && (SysEnv.IpfsRun == 0 || IpfsState == RUN) {
		return true
	}
	return false
}

func initEnv() error {
	er := cli.GetCoinbase()
	if er != nil {
		return er
	}
	er = getUsrs()
	if er != nil {
		return er
	}
	return nil
}

func getUsrs() error {
	// TODO:implements setting to global member
	ret, er := alg.UserMap_GetUsrs(ApEnv.UsrMap)
	if er != nil {
		fmt.Printf("exception occued in getUsers1:%s\n", er.Error())
		return er
	}
	fmt.Printf("AddressList:%s\n", strings.Join(ret, "\n"))
	Iwasmined := false
	for _, adrs := range ret {
		fmt.Printf("ards:%s\n", adrs)
		if adrs == cli.Coinbase {
			Iwasmined = true
		}
	}
	if !Iwasmined {
		// TODO:change to correct value
		fmt.Printf("ADDRESS:%s\n", ApEnv.UsrMap)
		tx, er := alg.UserMap_Reg(ApEnv.UsrMap, "0xaAaAfFfzfz", "for test")
		if er != nil {
			fmt.Printf("exception occued in getUsers2:%s\n", er.Error())
			return er
		}
		ret = append(ret, cli.Coinbase)
		go func() {
			for !Iwasmined {
				time.Sleep(time.Second * 10)
				adrs, cont, isnew, er := alg.UserMap_CheckReg(tx)
				if adrs != "" {
					fmt.Printf("I wasmined!:%s:%s:%s\n", adrs, cont, isnew)
					Iwasmined = true
				}
				if er != nil {
					fmt.Printf("%s", er.Error())
				}
			}
		}()
	}
	return nil
}
