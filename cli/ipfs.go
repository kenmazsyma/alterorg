// cli/ipfs.go

package cli

import (
	"errors"
	"fmt"
	ipfs "github.com/ipfs/go-ipfs-api"
	pb "github.com/ipfs/go-ipfs/unixfs/pb"
	proto "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/proto"
	"io"
	"os"
)

var ipfsurl string
var shell *ipfs.Shell
var myid string

func InitIpfs(url string) error {
	ipfsurl = url
	shell = ipfs.NewShell(ipfsurl)
	out, err := shell.ID()
	if err != nil {
		return err
	}
	myid = out.ID
	return nil
}

func IpfsAddFile(path string) (string, error) {
	file, er := os.Open(path)
	if er != nil {
		return "", er
	}
	return IpfsAdd(file)
}

func IpfsAdd(reader io.Reader) (string, error) {
	hash, er := shell.Add(reader)
	if er != nil {
		return "", er
	}
	return hash, nil
}

func IpfsGet(hash, path string) error {
	return shell.Get(hash, path)
}

func IpfsBlockGet(hash string) ([]byte, error) {
	return shell.BlockGet(hash)
}

var DIR_IPFS_BOARD = "boardforalterorg"
var dIR_IPFS = "/ipfs/"
var hASH_EMPTY_FILE = "QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH"
var hASH_EMPTY_DIR = "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"
var ERR_IpfsCreateIpfsDir_01 = "01" // ipns link is created as a file
var ERR_IpfsCreateIpfsDir_02 = "02" // dir which is used for boards of alterorg is created as a file
func IpfsCreateBoardDir() error {
	adrs, err := getIpnsAddr()
	if err != nil {
		return err
	}
	obj, err := shell.ObjectGet(adrs)
	if err != nil {
		return err
	}
	if len(obj.Links) == 0 && adrs != dIR_IPFS+hASH_EMPTY_DIR { // link to a file
		return errors.New(ERR_IpfsCreateIpfsDir_01)
	} else { // link to a directory
		lst, err := shell.List(adrs)
		if err != nil {
			return err
		}
		nodir := true
		for _, v := range lst {
			if v.Name == DIR_IPFS_BOARD {
				if v.Type == 2 { // 1:dir 2:file
					return errors.New(ERR_IpfsCreateIpfsDir_02)
				}
				nodir = false
				break
			}
		}
		// if the directory is not created yet
		if nodir {
			obj.Links = append(obj.Links, ipfs.ObjectLink{Name: DIR_IPFS_BOARD, Hash: hASH_EMPTY_DIR /*, Size: 3*/})
			size := uint64(0)
			buf, err := proto.Marshal(&pb.Data{Type: pb.Data_Directory.Enum(), Data: []byte(""), Filesize: &size})
			if err != nil {
				return err
			}
			obj.Data = string(buf)
			if err != nil {
				return err
			}
			newhash, err := shell.ObjectPut(obj)
			if err != nil {
				return err
			}
			fmt.Printf("NEW_NS:%s\n", newhash)
			err = shell.Publish("", newhash)
			if err != nil {
				return err
			}
		}
		fmt.Printf("Dir5\n")
	}
	return nil
}

func getIpnsAddr() (string, error) {
	adrs, err := shell.Resolve(myid)
	if err != nil {
		return "", err
	}
	return adrs, nil
}

var ERR_IpfsWriteToBoard_01 = "01" // dir for boards is not created yet
var ERR_IpfsWriteToBoard_02 = "02" // failed to get dir for boards
func IpfsWriteToBoard(data string, n string) error {
	size := uint64(len(data))
	buf, err := proto.Marshal(&pb.Data{Type: pb.Data_File.Enum(), Data: []byte(data), Filesize: &size})
	if err != nil {
		return err
	}
	objwr := ipfs.IpfsObject{Links: []ipfs.ObjectLink{}, Data: string(buf)}
	hash, err := shell.ObjectPut(&objwr)
	if err != nil {
		return err
	}
	fmt.Printf("BoardHash : %s\n", hash)
	nsadrs, err := getIpnsAddr()
	if err != nil {
		return err
	}
	// boarddir
	obj, err := shell.ObjectGet(nsadrs + "/" + DIR_IPFS_BOARD)
	if err != nil {
		if err.Error()[0:6] == "no link" {
			return errors.New(ERR_IpfsWriteToBoard_01)
		}
		return err
	}
	obj.Links = append(obj.Links, ipfs.ObjectLink{Name: n, Hash: hash, Size: uint64(len(data) + 100)})
	hash, err = shell.ObjectPut(obj)
	if err != nil {
		return err
	}
	fmt.Printf("NewHash for boardforalterorg : %s\n", hash)

	// ipnsdir
	nsobj, err := shell.ObjectGet(nsadrs)
	if err != nil {
		return err
	}
	found := false
	for i, v := range nsobj.Links {
		if v.Name == DIR_IPFS_BOARD {
			v.Hash = hash
			//	v.Size = 100
			nsobj.Links[i] = v
			found = true
			break
		}
	}
	if found == false {
		return errors.New(ERR_IpfsWriteToBoard_02)
	}
	nwnshash, err := shell.ObjectPut(nsobj)
	fmt.Printf("New Ns Hash including new boarddir:%s\n", nwnshash)
	// TODO:implements a process for async after implementing "wait" option to IPFS
	err = shell.Publish("", nwnshash)
	if err != nil {
		fmt.Printf("14\n")
		return err
	}

	return nil
}
