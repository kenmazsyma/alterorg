// cli/eth.go

package cli

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

func InitEth(url string) error {
	baseurl = url
	//	err := getCoinbase()
	//	if err != nil {
	//		return err
	//	}
	return nil
}

func GetCoinbase() error {
	arg := []Unknown{""}
	fmt.Printf("getcoinbase::::%s\n", baseurl)
	if er := Request(baseurl, "eth_coinbase", arg, &Coinbase); er != nil {
		fmt.Printf("\ngetCoinbase:%s\n", er.Error())
		return er
	}
	fmt.Printf("coinbase:%s\n", Coinbase)
	return nil
}

func GetEthListening() (bool, error) {
	arg := []Unknown{""}
	var ret bool
	fmt.Printf("GetEthListening::::%s\n", baseurl)
	if er := Request(baseurl, "net_listening", arg, &ret); er != nil {
		fmt.Printf("\nGetEthListening:%s\n", er.Error())
		return false, er
	}
	fmt.Printf("ret:%d", ret)
	return true, nil
}

func NewContract(code string, ab abi.ABI, param ...interface{}) (string, error) {
	enc, er := ab.Pack("", param...)
	if er != nil {
		return "", er
	}
	encst := hex.EncodeToString(enc)
	fmt.Printf("!!!code:%s\n", encst)
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
	fmt.Printf("Send:%s\n", name)
	fmt.Printf("Coinbase:%s\n", Coinbase)
	// TODO:set correct Gas value
	arg := []Unknown{argCall{From: Coinbase, To: to, Data: "0x" + encst, Gas: 2000000}}
	var data string
	if er := Request(baseurl, "eth_sendTransaction", arg, &data); er != nil {
		return "", er
	}
	fmt.Printf("Send2:%s\n", data)
	return data, nil
}

func Call(to string, ret interface{}, name string, ab abi.ABI, param ...interface{}) error {
	enc, er := ab.Pack(name, param...)
	if er != nil {
		fmt.Printf("Check point1\n")
		return er
	}
	encst := hex.EncodeToString(enc)
	fmt.Printf("Call:%s\n", name)
	// TODO:set right Gas value
	arg := []Unknown{argCall{From: Coinbase, To: to, Data: "0x" + encst, Gas: 20000000}, "latest"}
	var data string
	if er := Request(baseurl, "eth_call", arg, &data); er != nil {
		return er
	}
	fmt.Print("11111:" + data + "\n")
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
		fmt.Print(er.Error())
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
