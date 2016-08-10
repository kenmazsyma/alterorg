// alg/board.go

package alg

import (
	"../cli"
	"../cmn"
	"fmt"
	proto "github.com/gogo/protobuf/proto"
	ipfs "github.com/ipfs/go-ipfs-api"
	pb "github.com/ipfs/go-ipfs/unixfs/pb"
	"io/ioutil"
)

var dIR_IPFS = "/ipfs/"
var hASH_EMPTY_FILE = "QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH"
var hASH_EMPTY_DIR = "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"

type ErrCode string

const (
	ERR_IpfsCreateIpfsDir_01 ErrCode = "01" // ipns link is created as a file
	ERR_IpfsCreateIpfsDir_02 ErrCode = "02" // dir which is used for boards of alterorg is created as a file
)

func logBoard(txt string, args ...interface{}) {
	cmn.Log("[BOARD]", txt, args...)
}

func IpfsCreateBoardDir(name string) error {
	if err := cli.ChkNameResolved(); err != nil {
		return err
	}
	shell := cli.GetShell()
	obj, err := shell.ObjectGet(cli.GetIpnsAdrs())
	if err != nil {
		return err
	}
	if len(obj.Links) == 0 && cli.GetIpnsAdrs() != dIR_IPFS+hASH_EMPTY_DIR { // link to a file
		return makeError(ERR_IpfsCreateIpfsDir_01)
	} else { // link to a directory
		lst, err := shell.List(cli.GetIpnsAdrs())
		if err != nil {
			return err
		}
		nodir := true
		for _, v := range lst {
			if v.Name == name {
				if v.Type == 2 { // 1:dir 2:file
					return makeError(ERR_IpfsCreateIpfsDir_02)
				}
				nodir = false
				break
			}
		}
		// if the directory is not created yet
		if nodir {
			obj.Links = append(obj.Links, ipfs.ObjectLink{Name: name, Hash: hASH_EMPTY_DIR /*, Size: 3*/})
			size := uint64(0)
			buf, err := proto.Marshal(&pb.Data{Type: pb.Data_Directory.Enum(), Data: []byte(""), Filesize: &size})
			if err != nil {
				return err
			}
			obj.Data = string(buf)
			newhash, err := shell.ObjectPut(obj)
			if err != nil {
				return err
			}
			if err = cli.PutIpnsAdrs(newhash); err != nil {
				return err
			}
		}
	}
	return nil
}

const (
	ERR_IpfsWriteToBoard_01 ErrCode = "01" // dir for boards is not created yet
	ERR_IpfsWriteToBoard_02 ErrCode = "02" // failed to get dir for boards
)

func IpfsWriteToBoard(dir string, data string, n string) error {

	if err := cli.ChkStat(); err != nil {
		return err
	}
	size := uint64(len(data))
	buf, err := proto.Marshal(&pb.Data{Type: pb.Data_File.Enum(), Data: []byte(data), Filesize: &size})
	if err != nil {
		return err
	}
	shell := cli.GetShell()
	objwr := ipfs.IpfsObject{Links: []ipfs.ObjectLink{}, Data: string(buf)}
	hash, err := shell.ObjectPut(&objwr)
	if err != nil {
		return err
	}
	logBoard("BoardHash : %s", hash)
	// boarddir
	obj, err := shell.ObjectGet(cli.GetIpnsAdrs() + "/" + dir)
	if err != nil {
		if err.Error()[0:6] == "no link" {
			return makeError(ERR_IpfsWriteToBoard_01)
		}
		return err
	}
	obj.Links = append(obj.Links, ipfs.ObjectLink{Name: n, Hash: hash, Size: uint64(len(data) + 100)})
	hash, err = shell.ObjectPut(obj)
	if err != nil {
		return err
	}
	logBoard("NewHash for boardforalterorg : %s", hash)

	// ipnsdir
	nsobj, err := shell.ObjectGet(cli.GetIpnsAdrs())
	if err != nil {
		return err
	}
	found := false
	for i, v := range nsobj.Links {
		if v.Name == dir {
			v.Hash = hash
			nsobj.Links[i] = v
			found = true
			break
		}
	}
	if found == false {
		return makeError(ERR_IpfsWriteToBoard_02)
	}
	nwnshash, err := shell.ObjectPut(nsobj)
	if err = cli.PutIpnsAdrs(nwnshash); err != nil {
		return err
	}
	return nil
}

func IpfsListBoard(adrs string) ([][]string, error) {
	list, err := Assembly_GetParticipants(adrs)
	postlist := []*ipfs.LsLink{}
	if err != nil {
		return nil, err
	}
	ret := [][]string{}
	if err := cli.ChkStat(); err != nil {
		return nil, err
	}
	shell := cli.GetShell()
	for _, v := range list {
		user, err := UserMap_GetMappedUser(v)
		if err != nil {
			fmt.Printf("errored when founding user(%s) : %s\n", v, err.Error())
			continue
		}
		if user == "" {
			fmt.Printf("user(%s) is not found\n", v)
			continue
		}
		info, err := User_GetInfo(user)
		if err != nil {
			return nil, err
		}
		elms, err := shell.List("/ipns/" + info[1] + "/" + adrs)
		if err != nil {
			return nil, err
		}
		postlist = append(postlist, elms...)
	}
	// TODO:need to sort the list
	for _, lnk := range postlist {
		rc, err := shell.Cat(lnk.Hash)
		if err != nil {
			fmt.Printf(":%s(Hash:%s)\n", err.Error(), lnk.Hash)
			continue
		}
		buf, err := ioutil.ReadAll(rc)
		ret = append(ret, []string{lnk.Name, string(buf)})
	}
	return ret, nil
}
