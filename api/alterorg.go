package api

import (
	"../cmn"
	"fmt"
)

type ArgUpdSet struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type AlterOrg struct {
}

func NewAlterorg() *AlterOrg {
	return &AlterOrg{}
}

func (self *AlterOrg) QuerySetting(name []string, rslt *[]string) error {
	var er error
	*rslt, er = cmn.QueryEnv(name)
	if er != nil {
		fmt.Print(er.Error())
		return er
	}
	return nil
}

func (self *AlterOrg) UpdateSetting(val []ArgUpdSet, rslt *bool) error {
	*rslt = false
	for i := range val {
		er := cmn.UpdateEnv(val[i].Key, val[i].Val)
		if er != nil {
			fmt.Print(er.Error())
			return er
		}
	}
	er := cmn.SaveEnv("alterorg.json")
	if er != nil {
		fmt.Print(er.Error())
		return er
	}
	*rslt = true
	return nil
}

// register file to IPFS
func (self *AlterOrg) SaveFile(path string, rslt *string) error {
	return nil
}

// initialize ethereum, IPFS, etc
func (self *AlterOrg) Initialize(prm string, rslt *string) error {
	return nil
}

func (self *AlterOrg) Start(prm string, rslt *string) error {
	cmn.Start()
	return nil
}

func (self *AlterOrg) Stop(prm string, rslt *string) error {
	cmn.Stop()
	return nil
}
