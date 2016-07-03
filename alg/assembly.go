// alg/assembly.go

package alg

import (
	"../cli"
	"../cmn"
	sol "../solidity"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	//"strconv"
	"strings"
)

func logAssembly(txt string, args ...interface{}) {
	cmn.Log(LBL_ASSEMBLY, txt, args...)
}

type NotifyMine func(tx string, err error)

func Assembly_CheckMine(tx string) (string, error) {
	res, err := cli.CheckContractTransaction(tx)
	if err != nil {
		return "", err
	}
	if len(res.CA) > 0 {
		return res.CA, nil
	}
	return "", nil
}

func NewAssembly(name string) (string, error) {
	address, err := cli.NewContract(sol.Bin_Assembly, sol.Abi_Assembly, name)
	if err != nil {
		return "", err
	}
	return address, nil
}

func Assembly_ChkCreated(tx string) (string, error) {
	funcname := "onCreated"
	res, err := cli.CheckContractTransaction(tx)
	if err != nil {
		return "", err
	}
	if len(res.LOG) == 0 {
		return "", err
	}
	logAssembly("%s : number:%d, data:%s", funcname, len(res.LOG), res.LOG[0].Data)
	bdata, err := hex.DecodeString(res.LOG[0].Data[2:])
	if err != nil {
		return "", err
	}
	var (
		adrs1 = new(common.Address)
	)
	//ret := []interface{}{adrs1}
	if err = sol.Abi_Assembly.Unpack(adrs1, funcname, bdata); err != nil {
		return "", err
	}
	return adrs1.Hex(), nil
}

func Assembly_AddPerson(adrs string, list []string) (string, error) {
	funcname := "addPerson"
	for _, v := range list {
		if !checkAddress(v) {
			return "", errors.New("param for address is not correct format")
		}
	}
	param := [][]string{list}
	tx, err := cli.Send(adrs, funcname, sol.Abi_Assembly, param)
	if err != nil {
		return "", err
	}
	logAssembly("%s : %s", funcname, tx)
	return tx, nil
}

func Assembly_ChkAddedPerson(tx string) ([]string, error) {
	funcname := "onAddedPerson"
	res, err := cli.CheckContractTransaction(tx)
	if err != nil {
		return nil, err
	}
	if len(res.LOG) == 0 {
		return nil, err
	}
	logAssembly("%s : number:%d, data:%s", funcname, len(res.LOG), res.LOG[0].Data)
	bdata, err := hex.DecodeString(res.LOG[0].Data)
	if err != nil {
		return nil, err
	}
	var ret []string
	err = sol.Abi_Assembly.Unpack(ret, funcname, bdata)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func Assembly_GetName(adrs string) (string, error) {
	funcname := "getName"
	if !checkAddress(adrs) {
		return "", errors.New("param for address is not correct format")
	}
	var ret string
	err := cli.Call(adrs, ret, funcname, sol.Abi_Assembly)
	if err != nil {
		return "", err
	}
	logAssembly("%s : %s", funcname, ret)
	return ret, nil
}

func Assembly_GetBasicInfo(adrs string) ([]string, error) {
	funcname := "getBasicInfo"
	if !checkAddress(adrs) {
		return nil, errors.New("param for address is not correct format")
	}
	var (
		name     = new(string)
		proposal = new(string)
		arbiter  = new(common.Address)
		version  = /*new(uint)*/ big.NewInt(0)
	)
	ret := []interface{}{name, proposal, arbiter, &version}

	if err := cli.Call(adrs, &ret, funcname, sol.Abi_Assembly); err != nil {
		return nil, err
	}
	return []string{*name, *proposal, arbiter.Hex(), version.String()}, nil
}

func Assembly_GetProposal(adrs string) (string, string, string, error) {
	funcname := "getProposal"
	if !checkAddress(adrs) {
		return "", "", "", errors.New("param for address is not correct format")
	}
	var ret []string
	err := cli.Call(adrs, ret, funcname, sol.Abi_Assembly)
	if err != nil {
		return "", "", "", err
	}
	logAssembly("%s : %s", funcname, strings.Join(ret, ","))
	return ret[0], ret[1], ret[2], nil
}

func Assembly_RevisionProposal(adrs string, doc string, discuss string) (string, error) {
	funcname := "revisionProposal"
	if !checkAddress(adrs) {
		return "", errors.New("param for address is not correct format")
	}
	if !checkIPFSHash(doc) {
		return "", errors.New("param for proposal is not correct format")
	}
	if !checkIPFSHash(discuss) {
		return "", errors.New("param for discuss is not correct format")
	}
	param := []string{doc, discuss}
	tx, err := cli.Send(adrs, funcname, sol.Abi_Assembly, param)
	if err != nil {
		return "", err
	}
	logAssembly("%s : %s", funcname, tx)
	return tx, nil
}

type typeCheckRevision struct {
	Adrs string `json:"address"`
	Ver  uint   `json:"version"`
}

func Assembly_CheckRevision(tx string) (string, uint, error) {
	funcname := "onRevisionedProposal"
	res, err := cli.CheckContractTransaction(tx)
	if err != nil {
		return "", 0, err
	}
	if len(res.LOG) == 0 {
		return "", 0, err
	}
	logAssembly("%s : number:%d, data:%s", funcname, len(res.LOG), res.LOG[0].Data)
	bdata, err := hex.DecodeString(res.LOG[0].Data)
	if err != nil {
		return "", 0, err
	}
	var ret typeCheckRevision
	err = sol.Abi_Assembly.Unpack(ret, funcname, bdata)
	if err != nil {
		return "", 0, err
	}
	return ret.Adrs, ret.Ver, nil
}
