package cmn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type SysEnvSet struct {
	EthRun      int    `json:"eth_run"`
	EthCmd      string `json:"eth_cmd"`
	EthPrm      string `json:"eth_prm"`
	IpfsRun     int    `json:"ipfs_run"`
	IpfsCmd     string `json:"ipfs_cmd"`
	IpfsPrm     string `json:"ipfs_prm"`
	EthUrl      string `json:"eth_url"`
	IpfsUrl     string `json:"ipfs_url"`
	DownloadDir string `json:"download_dir"`
	ApEnvPath   string `json:"apenv_path"`
}

var SysEnv SysEnvSet

func LoadSysEnv(path string) error {
	fin, er := os.Open(path)
	if er != nil {
		fmt.Print("failure to open env file\n")
		return er
	}
	defer fin.Close()
	buf, er := ioutil.ReadAll(fin)
	if er != nil {
		fmt.Print("failure to read env file\n")
		return er
	}
	SysEnv = SysEnvSet{}
	er = json.Unmarshal(buf, &SysEnv)
	if er != nil {
		fmt.Printf("env file is bad format:%s\n", er.Error())
		return er
	}
	return nil
}

func SaveSysEnv(path string) error {
	data, er := json.Marshal(SysEnv)
	if er != nil {
		return er
	}
	ioutil.WriteFile(path, data, 0644)
	fin, er := os.Open(path)
	if er != nil {
		fmt.Print("failure to open env file\n")
		return er
	}
	defer fin.Close()
	buf, er := ioutil.ReadAll(fin)
	if er != nil {
		fmt.Print("failure to read env file\n")
		return er
	}
	env := SysEnvSet{}
	er = json.Unmarshal(buf, &env)
	if er != nil {
		fmt.Print("env file is bad format\n")
		return er
	}
	return nil
}

func QuerySysEnv(prm []string) ([]string, error) {
	ret := []string{}
	for i := range prm {
		switch prm[i] {
		case "eth_cmd":
			ret = append(ret, SysEnv.EthCmd)
		case "eth_prm":
			ret = append(ret, SysEnv.EthPrm)
		case "ipfs_cmd":
			ret = append(ret, SysEnv.IpfsCmd)
		case "ipfs_prm":
			ret = append(ret, SysEnv.IpfsPrm)
		case "eth_url":
			ret = append(ret, SysEnv.EthUrl)
		case "ipfs_url":
			ret = append(ret, SysEnv.IpfsUrl)
		case "download_dir":
			ret = append(ret, SysEnv.DownloadDir)
		default:
			ret = append(ret, "")
		}
	}
	return ret, nil
}

func UpdateSysEnv(key string, val string) error {
	// TODO:will be more smaaaaart!
	switch key {
	case "eth_cmd":
		SysEnv.EthCmd = val
	case "eth_prm":
		SysEnv.EthPrm = val
	case "ipfs_cmd":
		SysEnv.IpfsCmd = val
	case "ipfs_prm":
		SysEnv.IpfsPrm = val
	case "eth_url":
		SysEnv.EthUrl = val
	case "ipfs_url":
		SysEnv.IpfsUrl = val
	default:
		return errors.New("'" + key + "' is not supported")
	}
	return nil
}
