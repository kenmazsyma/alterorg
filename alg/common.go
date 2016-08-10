package alg

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const (
	LBL_ASSEMBLY string = "assembly"
)

// TODO:merge with cli/eth.go
type decodeproc func(string, uint) (string, uint, error)

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
	"uint256": func(val string, start uint) (string, uint, error) {
		ret, er := strconv.ParseInt(val[start:start+64], 16, 32)
		if er != nil {
			return "", 0, er
		}
		return fmt.Sprintf("%d", ret), 64, nil
	},
}

func checkIPFSHash(hash string) bool {
	if m, e := regexp.MatchString("^[0x]*[a-zA-Z0-9]{46}$", hash); !m || e != nil {
		return false
	}
	return true
}

func extractFunctionABI(abif interface{}, name string) (interface{}, error) {
	abi := abif.([]interface{})
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

func extractInputABI(abi interface{}, name string) ([]interface{}, error) {
	ex, er := extractFunctionABI(abi, name)
	if er != nil {
		fmt.Printf("failure to extract ABI for %s\n", name)
		return nil, er
	}
	return ex.(map[string]interface{})["inputs"].([]interface{}), nil
}

func extractOutputABI(abi interface{}, name string) ([]interface{}, error) {
	ex, er := extractFunctionABI(abi, name)
	if er != nil {
		fmt.Printf("failure to extract ABI for %s\n", name)
		return nil, er
	}
	return ex.(map[string]interface{})["outputs"].([]interface{}), nil
}

func binToMap(abi []interface{}, data string) (map[string]string, error) {
	st := 0
	ret := map[string]string{}
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
		ret[v["name"].(string)] = code
	}
	return ret, nil

}

func checkAddress(address string) bool {
	if m, e := regexp.MatchString("^[0x]*[a-zA-Z0-9]{40}$", address); !m || e != nil {
		return false
	}
	return true
}

func makeError(msg ErrCode) error {
	return errors.New(string(msg))
}
