// alg/assembly.go

package alg

import (
	"../cli"
	sol "../solidity"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type NotifyMine func(tx string, err error)

func Assembly_CheckMine(tx string) (string, error) {
	res, er := cli.CheckContractTransaction(tx)
	if er != nil {
		return "", er
	}
	if len(res.CA) > 0 {
		return res.CA, nil
	}
	return "", nil
}

func NewAssembly(proposal string, discuss string) (string, error) {
	// the value of param must be the hash of IPFS
	if !checkIPFSHash(proposal) {
		return "", errors.New("param for proposal is not correct format")
	}
	if !checkIPFSHash(discuss) {
		return "", errors.New("param for discuss is not correct format")
	}
	param := []string{proposal, discuss}
	address, er := cli.NewContract(sol.Bin_Assembly, sol.Abi_Assembly, param)
	if er != nil {
		return "", er
	}
	return address, nil
}

func Assembly_GetProposal(address string) (string, string, string, error) {
	funcname := "getProposal"
	if !checkAddress(address) {
		return "", "", "", errors.New("param for address is not correct format")
	}
	var ret []string
	er := cli.Call(address, ret, funcname, sol.Abi_Assembly)
	if er != nil {
		return "", "", "", er
	}
	fmt.Print(strings.Join(ret, ",") + "\n")
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
	tx, er := cli.Send(adrs, funcname, sol.Abi_Assembly, param)
	if er != nil {
		return "", er
	}
	fmt.Print(tx + "\n")
	return tx, nil
}

type typeCheckRevision struct {
	Adrs string `json:"address"`
	Ver  uint   `json:"version"`
}

func Assembly_CheckRevision(tx string) (string, uint, error) {
	funcname := "onRevisionedProposal"
	res, er := cli.CheckContractTransaction(tx)
	if er != nil {
		return "", 0, er
	}
	if len(res.LOG) == 0 {
		return "", 0, er
	}
	fmt.Printf("res.LOG:%d\n", len(res.LOG))
	fmt.Printf("res.LOG[0].Data:%s\n", res.LOG[0].Data)
	bdata, er := hex.DecodeString(res.LOG[0].Data)
	if er != nil {
		return "", 0, er
	}
	var ret typeCheckRevision
	er = sol.Abi_Assembly.Unpack(ret, funcname, bdata)
	if er != nil {
		return "", 0, er
	}
	return ret.Adrs, ret.Ver, nil
}
