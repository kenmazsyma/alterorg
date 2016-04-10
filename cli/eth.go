// cli/eth.go

package cli

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

type resTransactionReceipt struct {
	//	TH  string              `json:"transactionHash"`
	//	TI  string              `json:"transactionIndex"`
	//	BN  string              `json:"blockNumber"`
	//	BH  string              `json:"blockHash"`
	//	CGU string              `json:"cumulativeGasUsed"`
	//	GU  string              `json:"gasUsed"`
	CA string `json:"contractAddress"`
	//	LOG []map[string]string `json:"logs"`
}

type encodeproc func(string) (string, error)
type decodeproc func(string, uint) (string, uint, error)

var funcEncode = map[string]encodeproc{
	"bytes": func(val string) (string, error) {
		return "", nil // Dynamic length
	},
	"address": func(val string) (string, error) {
		if m, e := regexp.MatchString("^0x[0-9a-f]*$", val); !m || e != nil {
			return "", e
		}
		padst := string(bytes.Repeat([]byte{0x30}, 24))
		return val[2:] + padst, nil
	},
	"uint256": func(val string) (string, error) {
		return "", nil
	},
	"bool": func(val string) (string, error) {
		return fmt.Sprintf("%032x", val), nil
	},
}

var funcEncodeD = map[string]encodeproc{
	"bytes": func(val string) (string, error) {
		len := len(val)
		ret := fmt.Sprintf("%064x", len)
		code := []byte(val)
		v := fmt.Sprintf("%02x", code)
		pad := (32 - (len % 32)) * 2
		padst := string(bytes.Repeat([]byte{0x30}, pad))
		return ret + v + padst, nil
	},
}

var funcDecode = map[string]decodeproc{
	"bytes": func(val string, start uint) (string, uint, error) {
		ret := ""
		pos, er := strconv.ParseInt(val[start:start+64], 16, 32)
		if er != nil {
			return "", 0, er
		}
		pos *= 2
		len, er := strconv.ParseInt(val[pos:pos+64], 16, 32)
		if er != nil {
			return "", 0, er
		}
		for i := 0; i <= int(len)-2; i += 2 {
			cd, _ := strconv.ParseInt(val[int(pos)+i+64:int(pos)+i+66], 16, 16)
			ret += string(cd)
		}
		return ret, 64, nil
	},
	"address": func(val string, start uint) (string, uint, error) {
		return val[start+24 : start+64], 64, nil
	},
}

func InitEth(url string) error {
	baseurl = url
	err := getCoinbase()
	if err != nil {
		return err
	}
	return nil
}

func getCoinbase() error {
	arg := []Unknown{""}
	if er := Request(baseurl, "eth_coinbase", arg, &Coinbase); er != nil {
		return er
	}
	return nil
}

func NewContract(code string, param []string, abi interface{}) (string, error) {
	abii := abi.(map[string]interface{})["inputs"].([]interface{})
	enc, er := EncodeParam(param, abii)
	if er != nil {
		return "", er
	}
	// TODO:set right GAS value
	arg := []Unknown{argNewContract{From: Coinbase, Data: code + enc, Gas: 20000000}}
	var address string
	if er := Request(baseurl, "eth_sendTransaction", arg, &address); er != nil {
		return "", er
	}
	return address, nil
}

func CheckContractTransaction(ts string) (string, error) {
	arg := []Unknown{ts}
	res := resTransactionReceipt{}
	if er := Request(baseurl, "eth_getTransactionReceipt", arg, &res); er != nil {
		return "", er
	}
	return res.CA, nil
}

func Send(to string, name string, param []string, abi interface{}) (string, error) {
	abii := abi.(map[string]interface{})["inputs"].([]interface{})
	enc, er := EncodeParam(param, abii)
	if er != nil {
		return "", er
	}
	fmt.Printf("Send:%s\n", name)
	selector, er := getFunctionSelector(name, abii)
	if er != nil {
		return "", er
	}
	// TODO:set right Gas value
	arg := []Unknown{argCall{From: Coinbase, To: to, Data: "0x" + selector + enc, Gas: 20000000}}
	var data string
	if er := Request(baseurl, "eth_sendTransaction", arg, &data); er != nil {
		return "", er
	}
	fmt.Printf("Send2:%s\n", data)
	return data, nil
}

func Call(to string, name string, param []string, abi interface{}) ([]string, error) {
	abii := abi.(map[string]interface{})["inputs"].([]interface{})
	abio := abi.(map[string]interface{})["outputs"].([]interface{})
	enc, er := EncodeParam(param, abii)
	if er != nil {
		return nil, er
	}
	fmt.Printf("Call:%s\n", name)
	selector, er := getFunctionSelector(name, abii)
	if er != nil {
		return nil, er
	}
	// TODO:set right Gas value
	arg := []Unknown{argCall{From: Coinbase, To: to, Data: "0x" + selector + enc, Gas: 20000000}, "latest"}
	var data string
	if er := Request(baseurl, "eth_call", arg, &data); er != nil {
		return nil, er
	}
	fmt.Print("11111:" + data + "\n")
	ret, er := DecodeParam(data, abio)
	return ret, er
}

func ExtractFunctionAPI(abi []interface{}, name string) (interface{}, error) {
	for i := 0; i < len(abi); i++ {
		if name == "" {
			if abi[i].(map[string]interface{})["type"] == "constructor" {
				return abi[i], nil
			}
		} else {
			if abi[i].(map[string]interface{})["name"] == name {
				return abi[i], nil
			}
		}
	}
	return nil, errors.New("function '" + name + "' is not found")
}

func getFunctionSelector(name string, abi []interface{}) (string, error) {
	lst := []string{}
	for i := range abi {
		v := abi[i].(map[string]interface{})
		lst = append(lst, v["type"].(string))
	}
	val, err := Sha3ForString(name + "(" + strings.Join(lst, ",") + ")")
	if err != nil {
		return "", err
	}
	return val[2:10], nil
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

func EncodeParam(param []string, abi []interface{}) (string, error) {
	// fixed
	prm := []string{}
	pos := 0
	var l int
	if len(param) != len(abi) {
		return "",
			errors.New(fmt.Sprintf("invalid number of params. %d is corect, but %d",
				len(abi), len(param)))
	}
	for i := range abi {
		v := abi[i].(map[string]interface{})
		f := funcEncode[v["type"].(string)]
		if f == nil {
			return "", errors.New("don't support type:" + v["type"].(string))
		}
		fmt.Printf("%s:%s\n", v["type"].(string), param)
		code, er := f(param[i])
		if er != nil {
			return "", er
		}
		prm = append(prm, code)
		l = len(code) / 2
		if l == 0 {
			pos += 32
		} else {
			pos += l
		}
	}
	// dinamic
	for i := range abi {
		v := abi[i].(map[string]interface{})
		if prm[i] != "" {
			continue
		}
		f := funcEncodeD[v["type"].(string)]
		if f == nil {
			return "", errors.New("don't support type:" + v["type"].(string))
		}
		code, er := f(param[i])
		if er != nil {
			return "", er
		}
		prm[i] = fmt.Sprintf("%064x", pos)
		prm = append(prm, code)
		pos += len(code) / 2
	}
	return strings.Join(prm, ""), nil
}

func DecodeParam(data string, abi []interface{}) (ret []string, er error) {
	st := 0
	ret = []string{}
	data = data[2:]
	for i := range abi {
		v := abi[i].(map[string]interface{})
		f := funcDecode[v["type"].(string)]
		if f == nil {
			return nil, errors.New("don't support type:" + v["type"].(string))
		}
		code, len, er := f(data, uint(st))
		if er != nil {
			return nil, er
		}
		st += int(len)
		ret = append(ret, code)
	}
	return ret, nil
}

func decodeIPFSHash(val string) string {
	ret := ""
	for i := 0; i <= len(val)-2; i += 2 {
		cd, _ := strconv.ParseInt(val[i:i+2], 16, 16)
		ret += string(cd)
	}
	return ret
}

func encodeIPFSHash(val string) string {
	return fmt.Sprintf("%x", []byte(val))
}

func encodeAddress(val string) string {
	if m, e := regexp.MatchString("^0x[0-9a-f]*$", val); !m || e != nil {
		// exclude "0x"
		return val[2:]
	}
	return ""
}
