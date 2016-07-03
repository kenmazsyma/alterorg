// cli/eth.go

package cli

import (
	"../cmn"
	"./abi"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var baseurl string
var Coinbase string

type argNewContract struct {
	From string `json:"from"`
	Data string `json:"data"`
	Gas  int    `json:"gas"`
}

type argCall struct {
	From string `json:"from"`
	To   string `json:"to"`
	Data string `json:"data"`
	Gas  int    `json:"gas"`
}

type ResEventLog struct {
	Address     string   `json:"address"`
	BlockHash   string   `json:"blockHash"`
	BlockNumber string   `json:"blockNumber"`
	Data        string   `json:"data"`
	LogIndex    string   `json:"logIndex"`
	Topics      []string `json:"topics"`
	TxHash      string   `json:"transactionHash"`
	TxIndex     string   `json:"transactionIndex"`
}

type ResTransactionReceipt struct {
	//	TH  string              `json:"transactionHash"`
	//	TI  string              `json:"transactionIndex"`
	//	BN  string              `json:"blockNumber"`
	//	BH  string              `json:"blockHash"`
	//	CGU string              `json:"cumulativeGasUsed"`
	//	GU  string              `json:"gasUsed"`
	CA  string        `json:"contractAddress"`
	LOG []ResEventLog `json:"logs"`
}

/*
func InitEth(url string) error {
	baseurl = url
	//	err := getCoinbase()
	//	if err != nil {
	//		return err
	//	}
	return nil
}*/

var s_Eth Status

const (
	STTS_ETH_NOT_START        Status = 0
	STTS_ETH_WAIT_STARTING    Status = 20010
	STTS_ETH_WAIT_LISTENER    Status = 20020
	STTS_ETH_GETTING_COINBASE Status = 20030
	STTS_ETH_STARTED          Status = 21000
	STTS_ETH_FAILED           Status = 29999
)

var ethCmd *exec.Cmd

func logEth(txt string, args ...interface{}) {
	cmn.Log(LBL_ETH, txt, args...)
}

func StartEth() {
	s_Eth = STTS_ETH_NOT_START
	go func() {
		if cmn.SysEnv.EthRun != 0 {
			s_Eth = STTS_ETH_WAIT_STARTING
			prm := splitArgs(cmn.SysEnv.EthPrm)
			logEth("EthPrm:%s", strings.Join(prm, ":::"))
			ethCmd = exec.Command(cmn.SysEnv.EthCmd, prm...)
			ethCmd.Stdout = os.Stdout
			//ethCmd.Stdin = &ethInput
			if err := ethCmd.Start(); err != nil {
				logEth("Failed to start ethereum : %s", err.Error())
				s_Eth = STTS_ETH_FAILED
				return
			}
		}
		baseurl = cmn.SysEnv.EthUrl
		s_Eth = STTS_ETH_WAIT_LISTENER
		for i := 0; i < 10; i++ {
			time.Sleep(2 * time.Second)
			start, err := getEthListening()
			if err != nil {
				logEth("Waiting listener : %d", i+1)
			} else {
				if start {
					// TODO:now searching a way of getting unlock timing
					// if delete below sleep, "account is locked" error occurs
					// because unlock proc in geth is not called yet at this timing.
					time.Sleep(2 * time.Second)
					s_Eth = STTS_ETH_GETTING_COINBASE
					break
				}
			}
		}
		if s_Eth != STTS_ETH_GETTING_COINBASE {
			logEth("Failed to start listner for ehtereum")
			TermEth(STTS_ETH_FAILED)
			return
		}
		if err := getCoinbase(); err != nil {
			logEth("Failed to get coinbase : %s", err.Error())
			TermEth(STTS_ETH_FAILED)
			return
		}
		s_Eth = STTS_ETH_STARTED
	}()
}

func TermEth(stts Status) {
	logEth("Terminating ethereum...")
	if ethCmd != nil {
		ethCmd.Process.Signal(syscall.SIGINT)
		ethCmd = nil
	}
	s_Eth = stts
}

func GetEthStatus() Status {
	return s_Eth
}

func getCoinbase() error {
	arg := []Unknown{""}
	if er := Request(baseurl, "eth_coinbase", arg, &Coinbase); er != nil {
		logEth("getCoinbase:%s", er.Error())
		return er
	}
	logEth("coinbase:%s", Coinbase)
	return nil
}

func getEthListening() (bool, error) {
	arg := []Unknown{""}
	var ret bool
	if er := Request(baseurl, "net_listening", arg, &ret); er != nil {
		logEth("GetEthListening:%s", er.Error())
		return false, er
	}
	logEth("ret:%d", ret)
	return true, nil
}

func NewContract(code string, ab abi.ABI, param ...interface{}) (string, error) {
	enc, er := ab.Pack("", param...)
	if er != nil {
		return "", er
	}
	encst := hex.EncodeToString(enc)
	logEth("!!!code:%s", encst)
	// TODO:set correct GAS value
	arg := []Unknown{argNewContract{From: Coinbase, Data: code + encst, Gas: 2000000}}
	var address string
	if er := Request(baseurl, "eth_sendTransaction", arg, &address); er != nil {
		return "", er
	}
	return address, nil
}

func CheckContractTransaction(ts string) (ResTransactionReceipt, error) {
	arg := []Unknown{ts}
	res := ResTransactionReceipt{}
	if er := Request(baseurl, "eth_getTransactionReceipt", arg, &res); er != nil {
		return res, er
	}
	return res, nil
}

func Send(to string, name string, ab abi.ABI, param ...interface{}) (string, error) {
	enc, er := ab.Pack(name, param...)
	if er != nil {
		return "", er
	}
	encst := hex.EncodeToString(enc)
	logEth("Send:%s", name)
	logEth("Coinbase:%s", Coinbase)
	// TODO:set correct Gas value
	arg := []Unknown{argCall{From: Coinbase, To: to, Data: "0x" + encst, Gas: 2000000}}
	var data string
	if er := Request(baseurl, "eth_sendTransaction", arg, &data); er != nil {
		return "", er
	}
	logEth("Send2:%s", data)
	return data, nil
}

func Call(to string, ret interface{}, name string, ab abi.ABI, param ...interface{}) error {
	enc, er := ab.Pack(name, param...)
	if er != nil {
		logEth("Check point1")
		return er
	}
	encst := hex.EncodeToString(enc)
	logEth("Call:%s", name)
	// TODO:set right Gas value
	arg := []Unknown{argCall{From: Coinbase, To: to, Data: "0x" + encst, Gas: 20000000}, "latest"}
	var data string
	if er := Request(baseurl, "eth_call", arg, &data); er != nil {
		return er
	}
	logEth("Received Value[%s]:%s", name, data)
	bdata, er := hex.DecodeString(data[2:])
	if er != nil {
		return er
	}
	if len(bdata) == 0 {
		return nil
	}
	er = ab.Unpack(ret, name, bdata)
	return er
}

func Sha3(param string) (string, error) {
	arg := []Unknown{param}
	var data string
	if er := Request(baseurl, "web3_sha3", arg, &data); er != nil {
		logEth(er.Error())
		return "", er
	}
	return data, nil
}

func Sha3ForString(param string) (string, error) {
	return Sha3(StToHex(param))
}

func StToHex(val string) string {
	return "0x" + fmt.Sprintf("%x", []byte(val))
}
