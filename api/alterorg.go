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

type Alterorg struct {
}

func NewAlterorg() *Alterorg {
	return &Alterorg{}
}

func (self *Alterorg) QuerySetting(name []string, rslt *[]string) error {
	var er error
	*rslt, er = cmn.QuerySysEnv(name)
	if er != nil {
		fmt.Print(er.Error())
		return er
	}
	return nil
}

func (self *Alterorg) UpdateSetting(val []ArgUpdSet, rslt *bool) error {
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
func (self *Alterorg) SaveFile(path string, rslt *string) error {
	var er error
	*rslt, er = cli.IpfsAddFile(path)
	if er != nil {
		fmt.Printf("SaveFile:%s\n", er.Error())
		return er
	}
	return nil
}

// register file to IPFS
func (self *Alterorg) SaveData(data string, rslt *string) error {
	var er error
	reader := strings.NewReader(data)
	*rslt, er = cli.IpfsAdd(reader)
	if er != nil {
		fmt.Printf("SaveData:%s\n", er.Error())
		return er
	}
	return nil
}

func (self *Alterorg) GetFile(arg ArgGetFile, rslt *string) error {
	var er error
	er = cli.IpfsGet(arg.Hash, arg.Path)
	if er != nil {
		fmt.Printf("GetData:%s\n", er.Error())
		return er
	}
	return nil
}

func (self *Alterorg) GetData(hash string, rslt *string) error {
	data, er := cli.IpfsBlockGet(hash)
	if er != nil {
		fmt.Printf("GetData:%s\n", er.Error())
		return er
	}
	*rslt = string(data)
	return nil
}

// return system status
/*func (self *Alterorg) GetStatus(prm string, rslt *string) error {
	switch {
	case cmn.EthState == cmn.RUN && cmn.IpfsState == cmn.RUN:
		*rslt = "RUN"
	case cmn.EthState == cmn.ERROR || cmn.IpfsState == cmn.ERROR:
		*rslt = "ERROR"
	default:
		*rslt = "INIT"
	}
	return nil
}*/

func (self *Alterorg) QueryAssemblyLst(prm string, rslt *[]string) error {
	var er error
	*rslt, er = cmn.QueryAssemblyList()
	if er != nil {
		return er
	}
	return nil
}

// TODO:support an orgs isn't mined yet(save the tx hash)
func (self *Alterorg) UpdateAssemblyLst(prm []string, rslt *string) error {
	var er error
	fmt.Printf("UpdateAssemblyLst\n")
	er = cmn.UpdateAssemblyList(prm)
	if er != nil {
		return er
	}
	er = cmn.SaveApEnv(cmn.SysEnv.ApEnvPath)
	if er != nil {
		return er
	}
	return nil
}
