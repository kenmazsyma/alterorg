package api

import (
	"../cli"
	"../cmn"
	//"encoding/json"
	"fmt"
	"strings"
)

type ArgUpdSet struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type ArgGetFile struct {
	Hash string `json:"hash"`
	Path string `json:"path"`
}

type AlterOrg struct {
}

func NewAlterorg() *AlterOrg {
	return &AlterOrg{}
}

func (self *AlterOrg) QuerySetting(name []string, rslt *[]string) error {
	var er error
	*rslt, er = cmn.QuerySysEnv(name)
	if er != nil {
		fmt.Print(er.Error())
		return er
	}
	return nil
}

func (self *AlterOrg) UpdateSetting(val []ArgUpdSet, rslt *bool) error {
	*rslt = false
	for i := range val {
		er := cmn.UpdateSysEnv(val[i].Key, val[i].Val)
		if er != nil {
			fmt.Print(er.Error())
			return er
		}
	}
	er := cmn.SaveSysEnv("alterorg.json")
	if er != nil {
		fmt.Print(er.Error())
		return er
	}
	*rslt = true
	return nil
}

// register file to IPFS
func (self *AlterOrg) SaveFile(path string, rslt *string) error {
	var er error
	*rslt, er = cli.IpfsAddFile(path)
	if er != nil {
		fmt.Printf("SaveFile:%s\n", er.Error())
		return er
	}
	return nil
}

// register file to IPFS
func (self *AlterOrg) SaveData(data string, rslt *string) error {
	var er error
	reader := strings.NewReader(data)
	*rslt, er = cli.IpfsAdd(reader)
	if er != nil {
		fmt.Printf("SaveData:%s\n", er.Error())
		return er
	}
	return nil
}

func (self *AlterOrg) GetFile(arg ArgGetFile, rslt *string) error {
	var er error
	er = cli.IpfsGet(arg.Hash, arg.Path)
	if er != nil {
		fmt.Printf("GetData:%s\n", er.Error())
		return er
	}
	return nil
}

func (self *AlterOrg) GetData(hash string, rslt *string) error {
	data, er := cli.IpfsBlockGet(hash)
	if er != nil {
		fmt.Printf("GetData:%s\n", er.Error())
		return er
	}
	*rslt = string(data)
	return nil
}

// return system status
func (self *AlterOrg) GetStatus(prm string, rslt *string) error {
	switch {
	case cmn.EthState == cmn.RUN && cmn.IpfsState == cmn.RUN:
		*rslt = "RUN"
	case cmn.EthState == cmn.ERROR || cmn.IpfsState == cmn.ERROR:
		*rslt = "ERROR"
	default:
		*rslt = "INIT"
	}
	return nil
}

func (self *AlterOrg) QueryOrgLst(prm string, rslt *[]string) error {
	var er error
	*rslt, er = cmn.QueryOrgList()
	if er != nil {
		return er
	}
	return nil
}

func (self *AlterOrg) UpdateOrgLgst(prm []string, rslt *string) error {
	var er error
	er = cmn.UpdateOrgList(prm)
	if er != nil {
		return er
	}
	return nil
}
