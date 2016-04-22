package cmn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type EnvVal struct {
	EthCmd  string `json:"eth_cmd"`
	EthPrm  string `json:"eth_prm"`
	IpfsCmd string `json:"ipfs_cmd"`
	IpfsPrm string `json:"ipfs_prm"`
	EthUrl  string `json:"eth_url"`
	IpfsUrl string `json:"ipfs_url"`
}

var Env EnvVal

func LoadEnv(path string) error {
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
	Env = EnvVal{}
	er = json.Unmarshal(buf, &Env)
	if er != nil {
		fmt.Printf("env file is bad format:%s\n", er.Error())
		return er
	}
	return nil
}

func SaveEnv(path string) error {
	data, er := json.Marshal(Env)
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
	env := EnvVal{}
	er = json.Unmarshal(buf, &env)
	if er != nil {
		fmt.Print("env file is bad format\n")
		return er
	}
	return nil
}

func QueryEnv(prm []string) ([]string, error) {
	ret := []string{}
	for i := range prm {
		switch prm[i] {
		case "eth_cmd":
			ret = append(ret, Env.EthCmd)
		case "eth_prm":
			ret = append(ret, Env.EthPrm)
		case "ipfs_cmd":
			ret = append(ret, Env.IpfsCmd)
		case "ipfs_prm":
			ret = append(ret, Env.IpfsPrm)
		case "eth_url":
			ret = append(ret, Env.EthUrl)
		case "ipfs_url":
			ret = append(ret, Env.IpfsUrl)
		default:
			ret = append(ret, "")
		}
	}
	return ret, nil
}

func UpdateEnv(key string, val string) error {
	// TODO:will be more smaaaaart!
	switch key {
	case "eth_cmd":
		Env.EthCmd = val
	case "eth_prm":
		Env.EthPrm = val
	case "ipfs_cmd":
		Env.IpfsCmd = val
	case "ipfs_prm":
		Env.IpfsPrm = val
	case "eth_url":
		Env.EthUrl = val
	case "ipfs_url":
		Env.IpfsUrl = val
	default:
		return errors.New("'" + key + "' is not supported")
	}
	return nil
}
