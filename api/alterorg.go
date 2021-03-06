package api

import (
	"../alg"
	"../cli"
	"../cmn"
	//"encoding/json"
	"errors"
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
	var err error
	*rslt, err = cmn.QuerySysEnv(name)
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	return nil
}

func (self *Alterorg) UpdateSetting(val []ArgUpdSet, rslt *bool) error {
	*rslt = false
	for i := range val {
		err := cmn.UpdateSysEnv(val[i].Key, val[i].Val)
		if err != nil {
			fmt.Print(err.Error())
			return err
		}
	}
	err := cmn.SaveSysEnv("alterorg.json")
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	*rslt = true
	return nil
}

// register file to IPFS
func (self *Alterorg) SaveFile(path string, rslt *string) error {
	var err error
	*rslt, err = cli.IpfsAddFile(path)
	if err != nil {
		fmt.Printf("SaveFile:%s\n", err.Error())
		return err
	}
	return nil
}

// register file to IPFS
func (self *Alterorg) SaveData(data string, rslt *string) error {
	var err error
	reader := strings.NewReader(data)
	*rslt, err = cli.IpfsAdd(reader)
	if err != nil {
		fmt.Printf("SaveData:%s\n", err.Error())
		return err
	}
	return nil
}

func (self *Alterorg) GetFile(arg ArgGetFile, rslt *string) error {
	var err error
	err = cli.IpfsGet(arg.Hash, arg.Path)
	if err != nil {
		fmt.Printf("GetData:%s\n", err.Error())
		return err
	}
	return nil
}

func (self *Alterorg) GetData(hash string, rslt *string) error {
	data, err := cli.IpfsBlockGet(hash)
	if err != nil {
		fmt.Printf("GetData:%s\n", err.Error())
		return err
	}
	*rslt = string(data)
	return nil
}

func (self *Alterorg) GetEthStatus(prm string, rslt *string) error {
	stts := cli.GetEthStatus()
	switch {
	case stts == cli.STTS_ETH_STARTED:
		*rslt = "RUN"
	case stts == cli.STTS_ETH_FAILED:
		*rslt = "ERROR"
	default:
		*rslt = "INIT"
	}
	return nil
}
func (self *Alterorg) GetIpfsStatus(prm string, rslt *string) error {
	stts := cli.GetIpfsStatus()
	switch {
	case stts == cli.STTS_IPFS_STARTED:
		*rslt = "RUN"
	case stts == cli.STTS_IPFS_FAILED:
		*rslt = "ERROR"
	default:
		*rslt = "INIT"
	}
	return nil
}

func (self *Alterorg) QueryAssemblyLst(prm string, rslt *[]string) error {
	*rslt = cmn.QueryAssemblyList()
	return nil
}

// TODO:support an orgs isn't mined yet(save the tx hash)
func (self *Alterorg) UpdateAssemblyLst(prm []string, rslt *string) error {
	var err error
	fmt.Printf("UpdateAssemblyLst\n")
	cmn.UpdateAssemblyList(prm)
	err = cmn.SaveApEnv(cmn.SysEnv.ApEnvPath)
	if err != nil {
		return err
	}
	return nil
}

func (selct *Alterorg) AppendAssembly(adrs string, rslt *string) error {
	lst := append(cmn.QueryAssemblyList(), adrs)
	cmn.UpdateAssemblyList(lst)
	err := cmn.SaveApEnv(cmn.SysEnv.ApEnvPath)
	if err != nil {
		return err
	}
	return nil
}

func (self *Alterorg) WriteToBoard(prm []string, rslt *string) error {
	fmt.Printf("[Alterorg]%s, %s, %s\n", prm[0], prm[1], prm[2])

	if len(prm) != 3 {
		return errors.New("Invalid parameters")
	}
	return alg.IpfsWriteToBoard(prm[0], prm[1], prm[2])
}

func (self *Alterorg) ListBoard(prm string, rslt *[][]string) error {
	var err error
	*rslt, err = alg.IpfsListBoard(prm)
	if err != nil {
		return err
	}
	return nil
}

func (self *Alterorg) PrepareBoard(prm string, rslt *string) error {
	if err := alg.IpfsCreateBoardDir(prm); err != nil {
		return err
	}
	return nil
}
func (self *Alterorg) GetEthAddress(prm string, rslt *string) error {
	*rslt = cli.Coinbase
	return nil
}
