// alg/user.go

package alg

import (
	"../cli"
	"../cmn"
	sol "../solidity"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
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

func UserMap_RegOwn(address string, node string, name string) (string, error) {
	funcname := "reg"
	tx, err := cli.Send(address, funcname, sol.Abi_UserMap, []byte(node), name)
	if err != nil {
		return "", err
	}
	fmt.Print(tx + "\n")
	return tx, nil
}

type Status int

const (
	STTS_USER_NOT_GET    Status = 0
	STTS_USER_WAIT_ETH   Status = 10010
	STTS_USER_GETTING    Status = 10020
	STTS_USER_WAIT_REG   Status = 10030
	STTS_USER_REGISTERED Status = 10100
	STTS_USER_FAILED     Status = 19999
)

func logUser(txt string, args ...interface{}) {
	cmn.Log("user", txt, args...)
}

var s_User Status

func UserMap_Prepare() {
	go func() {
		s_User = STTS_USER_WAIT_ETH
		logUser("Wainting Ethereum & IPFS")
		for true {
			time.Sleep(1 * time.Second)
			ipfsstat := cli.GetIpfsStatus()
			if ipfsstat != cli.STTS_IPFS_STARTED && ipfsstat != cli.STTS_IPFS_RESOLVING_NAME {
				continue
			}
			if cli.GetEthStatus() != cli.STTS_ETH_STARTED {
				continue
			}
			logUser("Confirmed IPFS & Ethereum started")
			break
		}
		s_User = STTS_USER_GETTING
		UsrLst, err := UserMap_GetUsrs()
		if err != nil {
			logUser("Failed to get UserList:%s", err.Error())
			s_User = STTS_USER_FAILED
			return
		}
		logUser("AddressList:%s", strings.Join(UsrLst, "\n"))
		mined := false
		for _, adrs := range UsrLst {
			fmt.Printf("check list:%s, %s\n", adrs, cli.Coinbase)
			if adrs == cli.Coinbase {
				mined = true
			}
		}
		if !mined {
			s_User = STTS_USER_WAIT_REG
			// TODO:change to correct value
			tx, err := UserMap_RegOwn(cmn.ApEnv.UsrMap, cli.GetIpfsId(), "")
			if err != nil {
				logUser("Failed to regist my coount to UsrList:%s", err.Error())
				s_User = STTS_USER_FAILED
				return
			}
			UsrLst = append(UsrLst, cli.Coinbase)
			time.Sleep(time.Second * 3)
			for !mined {
				adrs, cont, isnew, err := UserMap_CheckRegOwn(tx)
				if adrs != "" {
					logUser("I was mined!:%s:%s:%s", adrs, cont, isnew)
					mined = true
					s_User = STTS_USER_REGISTERED
					break
				}
				if err != nil {
					logUser("%s", err.Error())
					s_User = STTS_USER_FAILED
					return
				}
				time.Sleep(time.Second * 1)
			}
		}
		s_User = STTS_USER_REGISTERED
		logUser("Success to get user lists")
	}()

}

func UserMap_CheckRegOwn(tx string) (string, string, bool, error) {
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

func UserMap_GetUsrs() ([]string, error) {
	funcname := "getAddresses"
	//param := []string{}
	var adss []ethcmn.Address
	err := cli.Call(cmn.ApEnv.UsrMap, &adss, funcname, sol.Abi_UserMap)
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, v := range adss {
		logUser("elm : %s", v.Hex())
		ret = append(ret, v.Hex())
	}
	return ret, nil
}

func User_GetInfo(address string) ([]string, error) {
	funcname := "getInfo"
	if checkAddress(address) == false {
		return nil, errors.New("param for address is not correct format")
	}
	var (
		ret1 = new(common.Address)
		ret2 = []byte{}
		ret3 = new(string)
	)
	ret := []interface{}{ret1, &ret2, ret3}
	if err := cli.Call(address, &ret, funcname, sol.Abi_User); err != nil {
		return nil, err
	}
	return []string{ret1.Hex(), string(ret2), *ret3}, nil
}

func UserMap_GetMappedUser(adrs4usr string) (string, error) {
	funcname := "getUser"

	if checkAddress(adrs4usr) == false {
		return "", errors.New("param for address of user is not correct format")
	}
	ret := common.Address{}
	prm := common.HexToAddress(adrs4usr)
	if err := cli.Call(cmn.ApEnv.UsrMap, &ret, funcname, sol.Abi_UserMap, prm); err != nil {
		return "", err
	}
	rethex := ret.Hex()
	if rethex == "0x0000000000000000000000000000000000000000" {
		return "", nil
	}
	return rethex, nil
}

func UserMap_Reg(node string, name string) (string, error) {
	funcname := "reg"
	//param := []string{node, name}
	tx, err := cli.Send(cmn.ApEnv.UsrMap, funcname, sol.Abi_UserMap, []byte(node), name)
	if err != nil {
		return "", err
	}
	return tx, nil
}

func UserMap_CheckReg(tx string) (string, error) {
	funcname := "onReg"
	res, err := cli.CheckContractTransaction(tx)
	if err != nil {
		return "", err
	}
	if len(res.LOG) == 0 {
		return "", err
	}
	bdata, err := hex.DecodeString(res.LOG[0].Data[2:])
	if err != nil {
		return "", err
	}
	var (
		adrs  = new(common.Address)
		con   = new(common.Address)
		isnew = new(bool)
	)
	ret := []interface{}{adrs, con, isnew}
	err = sol.Abi_UserMap.Unpack(&ret, funcname, bdata)
	if err != nil {
		return "", err
	}
	fmt.Printf("%s, %s, %s", adrs.Hex(), con.Hex(), *isnew)
	return con.Hex(), nil
}
