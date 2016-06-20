// alg/user.go

package alg

import (
	"../cli"
	"../cmn"
	sol "../solidity"
	"encoding/hex"
	"errors"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"runtime/debug"
	"strings"
	"time"
)

var UseLst []string

// UserMap is singlton contract. This function is not called in actually.
func NewUserMap() (string, error) {
	address, err := cli.NewContract(sol.Bin_UserMap, sol.Abi_UserMap)
	if err != nil {
		return "", err
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
	tx, err := cli.Send(address, funcname, sol.Abi_UserMap /*ethcmn.StringToAddress(adrs)*/, []byte(adrs), name)
	if err != nil {
		return "", err
	}
	fmt.Print(tx + "\n")
	return tx, nil
}

//type typeCheckReg struct {
//	Adrs  ethcmn.Address `json:"adrs"`
//	Cont  ethcmn.Address `json:"cont"`
//	IsNew bool           `json:"isNew"`
//}

type Status int

const (
	STTS_USER_NOT_GET    Status = 0
	STTS_USER_WAIT_ETH   Status = 10010
	STTS_USER_GETTING    Status = 10020
	STTS_USER_WAIT_REG   Status = 10030
	STTS_USER_REGISTERED Status = 10100
	STTS_USER_FAILED     Status = 19999
)

var s_User chan Status

func UserMap_Prepare() {
	s_User = make(chan Status, STTS_USER_NOT_GET)
	go func() {
		s_User <- STTS_USER_WAIT_ETH
		fmt.Println("Wainting Ethereum & IPFS")
		for true {
			time.Sleep(1 * time.Second)
			if cli.GetIpfsStatus() != cli.STTS_IPFS_STARTED {
				continue
			}
			if cli.GetEthStatus() != cli.STTS_ETH_STARTED {
				continue
			}
			break
		}
		s_User <- STTS_USER_GETTING
		UsrLst, err := UserMap_GetUsrs(cmn.ApEnv.UsrMap)
		if err != nil {
			fmt.Printf("Failed to get UserList:%s\n", err.Error())
			s_User <- STTS_USER_FAILED
			return
		}
		fmt.Printf("AddressList:%s\n", strings.Join(UsrLst, "\n"))
		mined := false
		for _, adrs := range UsrLst {
			if adrs == cli.Coinbase {
				mined = true
			}
		}
		if !mined {
			s_User <- STTS_USER_WAIT_REG
			// TODO:change to correct value
			tx, err := UserMap_Reg(cmn.ApEnv.UsrMap, "0xaAaAfFfzfz", "for test")
			if err != nil {
				fmt.Printf("Failed to regist my coount to UsrList:%s\n", err.Error())
				s_User <- STTS_USER_FAILED
				return
			}
			UsrLst = append(UsrLst, cli.Coinbase)
			time.Sleep(time.Second * 3)
			for !mined {
				adrs, cont, isnew, err := UserMap_CheckReg(tx)
				if adrs != "" {
					fmt.Printf("I was mined!:%s:%s:%s\n", adrs, cont, isnew)
					mined = true
					s_User <- STTS_USER_REGISTERED
					break
				}
				if err != nil {
					fmt.Printf("%s", err.Error())
					s_User <- STTS_USER_FAILED
					return
				}
				time.Sleep(time.Second * 1)
			}
		}
		s_User <- STTS_USER_REGISTERED
	}()

}

func UserMap_CheckReg(tx string) (string, string, bool, error) {
	funcname := "onReg"
	res, err := cli.CheckContractTransaction(tx)
	if err != nil {
		return "", "", false, err
	}
	if len(res.LOG) == 0 {
		return "", "", false, err
	}
	fmt.Printf("res.LOG:%d\n", len(res.LOG))
	fmt.Printf("res.LOG[0].Data:%s\n", res.LOG[0].Data)
	bdata, err := hex.DecodeString(res.LOG[0].Data[2:])
	if err != nil {
		return "", "", false, err
	}
	//var ret typeCheckReg
	var (
		var1 = new(ethcmn.Address)
		var2 = new(ethcmn.Address)
		var3 = new(bool)
	)
	//ret := []interface{}{new(ethcmn.Address), new(ethcmn.Address), new(bool)}
	ret := []interface{}{var1, var2, var3}
	err = sol.Abi_UserMap.Unpack(&ret, funcname, bdata)
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
	err := cli.Call(address, &adss, funcname, sol.Abi_UserMap)
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, v := range adss {
		ret = append(ret, v.Hex())
	}
	return ret, nil
}
