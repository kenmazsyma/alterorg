package api

import (
	"../alg"
	"../cli"
	"../cmn"
	"fmt"
)

type User struct {
}

type RsltUserGetInfo struct {
	Name      string `json:"name"`
	Adrs4Eth  string `json:"adrs4eth"`
	Adrs4Ipfs string `json:"adrs4ipfs"`
}

func NewUser() *User {
	return &User{}
}

func (self *User) GetInfo(adrs string, rslt *RsltUserGetInfo) error {
	info, err := alg.User_GetInfo(adrs)
	if err != nil {
		fmt.Print("%s:\n", err.Error())
		return err
	}
	*rslt = RsltUserGetInfo{Name: info[2], Adrs4Eth: info[0], Adrs4Ipfs: info[1]}
	return nil
}

func (self *User) GetMappedUser(adrs string, rslt *string) error {
	fmt.Printf("UserMap:%s\n", cmn.ApEnv.UsrMap)
	ret, err := alg.User_GetMappedUser(cmn.ApEnv.UsrMap, adrs)
	if err != nil {
		fmt.Print("%s:\n", err.Error())
		return err
	}
	*rslt = ret
	return nil
}

type ArgUserReg struct {
	Node string `json:"node"`
	Name string `json:"name"`
}

func (self *User) Reg(prm ArgUserReg, rslt *string) error {
	// TODO:change to get data from param
	ipfsid := cli.GetIpfsId()
	tx, err := alg.User_Reg(cmn.ApEnv.UsrMap, ipfsid, prm.Name)
	if err != nil {
		fmt.Printf("%s:\n", err.Error())
		return err
	}
	*rslt = tx
	return nil
}

func (self *User) CheckReg(tx string, rslt *string) error {
	adrs, err := alg.User_CheckReg(tx)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return err
	}
	*rslt = adrs
	return nil
}
