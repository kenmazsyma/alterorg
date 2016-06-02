// alg/user.go

package alg

import (
	"../cli"
	sol "../solidity"
	"encoding/hex"
	"errors"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"runtime/debug"
)

// UserMap is singlton contract. This function is not called in actually.
func NewUserMap() (string, error) {
	address, er := cli.NewContract(sol.Bin_UserMap, sol.Abi_UserMap)
	if er != nil {
		return "", er
	}
	debug.PrintStack()
	return address, nil
}

type typeReg struct {
	Node ethcmn.Address `json:"node"`
	Name string         `json:"name"`
}

func UserMap_Reg(address string, node string, name string) (string, error) {
	//param := typeReg{Node: ethcmn.StringToAddress(adrs), Name: name}
	funcname := "reg"
	//adrs:=Coinbase
	adrs := "0x1111111111222222222233333333334444444444"
	tx, er := cli.Send(address, funcname, sol.Abi_UserMap /*ethcmn.StringToAddress(adrs)*/, []byte(adrs), name)
	if er != nil {
		return "", er
	}
	fmt.Print(tx + "\n")
	return tx, nil
}

//type typeCheckReg struct {
//	Adrs  ethcmn.Address `json:"adrs"`
//	Cont  ethcmn.Address `json:"cont"`
//	IsNew bool           `json:"isNew"`
//}

func UserMap_CheckReg(tx string) (string, string, bool, error) {
	funcname := "onReg"
	res, er := cli.CheckContractTransaction(tx)
	if er != nil {
		return "", "", false, er
	}
	if len(res.LOG) == 0 {
		return "", "", false, er
	}
	fmt.Printf("res.LOG:%d\n", len(res.LOG))
	fmt.Printf("res.LOG[0].Data:%s\n", res.LOG[0].Data)
	bdata, er := hex.DecodeString(res.LOG[0].Data[2:])
	if er != nil {
		return "", "", false, er
	}
	//var ret typeCheckReg
	var (
		var1 = new(ethcmn.Address)
		var2 = new(ethcmn.Address)
		var3 = new(bool)
	)
	//ret := []interface{}{new(ethcmn.Address), new(ethcmn.Address), new(bool)}
	ret := []interface{}{var1, var2, var3}
	er = sol.Abi_UserMap.Unpack(&ret, funcname, bdata)
	//return ret[0].(*ethcmn.Address).Str(), ret[1].(*ethcmn.Address).Str(), *(ret[2].(*bool)), nil
	return var1.Hex(), var2.Hex(), *var3, nil
}

func UserMap_GetUsrs(address string) ([]string, error) {
	funcname := "getAddresses"
	if !checkAddress(address) {
		return nil, errors.New("param for address is not correct format")
	}
	//param := []string{}
	var adss []ethcmn.Address
	er := cli.Call(address, &adss, funcname, sol.Abi_UserMap)
	if er != nil {
		return nil, er
	}
	var ret []string
	for _, v := range adss {
		ret = append(ret, v.Hex())
	}
	return ret, nil
}
