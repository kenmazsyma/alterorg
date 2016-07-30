// cli/ipfs.go

package cli

import (
	"../cmn"
	"errors"
	"fmt"
	ipfs "github.com/ipfs/go-ipfs-api"
	pb "github.com/ipfs/go-ipfs/unixfs/pb"
	//proto "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/proto"
	proto "github.com/gogo/protobuf/proto"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"time"
)

const (
	STTS_IPFS_NOT_START       Status = 0
	STTS_IPFS_INITIALIZING    Status = 10010
	STTS_IPFS_STARTING        Status = 10020
	STTS_IPFS_GETTING_SYSINFO Status = 10030
	STTS_IPFS_RESOLVING_NAME  Status = 10040
	STTS_IPFS_STARTED         Status = 11000
	STTS_IPFS_FAILED          Status = 19999
)

var s_Ipfs Status
var ipfsurl string
var shell *ipfs.Shell
var myid string
var nsAdrs string
var ipfsCmd *exec.Cmd
var ipfsid string

func GetIpfsId() string {
	return ipfsid
}

func logIpfs(txt string, args ...interface{}) {
	cmn.Log(LBL_IPFS, txt, args...)
}

func GetIpfsStatus() Status {
	return s_Ipfs
}

func StartIpfs() {
	s_Ipfs = STTS_IPFS_NOT_START
	go func() {
		if cmn.SysEnv.IpfsRun != 0 {
			const (
				NOREPO string = "Error: no ipfs repo found"
			)
			out, err := exec.Command(cmn.SysEnv.IpfsCmd, "diag", "net").CombinedOutput()
			if err != nil {
				if string(out[0:25]) == NOREPO {
					s_Ipfs = STTS_IPFS_INITIALIZING
					logIpfs("IPFS is not initialized.\nNow initializing...")
					out, err = exec.Command(cmn.SysEnv.IpfsCmd, "init").Output()
					if err != nil {
						s_Ipfs = STTS_IPFS_FAILED
						logIpfs("Failed to execute ipfs init:%s", err.Error())
						return
					}
					logIpfs("Initialize IPFS:%s", out)
				} else {
					// It's OK because ipfs daemon is not run
				}
			}

			s_Ipfs = STTS_IPFS_STARTING
			logIpfs("Starting IPFS...")
			prm := splitArgs(cmn.SysEnv.IpfsPrm)
			ipfsCmd = exec.Command(cmn.SysEnv.IpfsCmd, prm...)
			if err = ipfsCmd.Start(); err != nil {
				s_Ipfs = STTS_IPFS_FAILED
				logIpfs("error in ipfs daemon:%s", err.Error())
				return
			}
		}
		var out *ipfs.IdOutput
		var err error
		for true {
			s_Ipfs = STTS_IPFS_GETTING_SYSINFO
			time.Sleep(1 * time.Second)
			logIpfs("Getting information from IPFS....")
			shell = ipfs.NewShell(cmn.SysEnv.IpfsUrl)
			out, err = shell.ID()
			if err != nil {
				//	logIpfs("error in shell.ID:%s", err.Error())
				//	TermIpfs(STTS_IPFS_FAILED)
				//	return
				continue
			}
			ipfsid = out.ID
			break
		}
		for true {
			s_Ipfs = STTS_IPFS_RESOLVING_NAME
			time.Sleep(1 * time.Second)
			logIpfs("Getting my ipns address")
			myid = out.ID
			if err := getIpnsAdrs(); err != nil {
				//	logIpfs("Failed to get IPNS address:%s", err.Error())
				//	TermIpfs(STTS_IPFS_FAILED)
				//	return
				continue
			}
			break
		}
		s_Ipfs = STTS_IPFS_STARTED
		//		if err := IpfsCreateBoardDir(); err != nil {
		//			logIpfs("Failed to create board dir:%s", err.Error())
		//			TermIpfs(STTS_IPFS_FAILED)
		//			return
		//		}
	}()
}

func TermIpfs(stts Status) {
	logIpfs("Terminationg Ipfs...")
	if ipfsCmd != nil {
		ipfsCmd.Process.Signal(syscall.SIGINT)
		ipfsCmd = nil
	}
	s_Ipfs = stts
}

func chkStat() error {
	if s_Ipfs != STTS_IPFS_STARTED && s_Ipfs != STTS_IPFS_RESOLVING_NAME {
		return errors.New("Ipfs is not started")
	}
	return nil
}

func chkNameResolved() error {
	if s_Ipfs != STTS_IPFS_STARTED {
		return errors.New("Ipfs is not started")
	}
	return nil
}

/*
func InitIpfs(url string) error {
	ipfsurl = url
	shell = ipfs.NewShell(ipfsurl)
	out, err := shell.ID()
	if err != nil {
		return err
	}
	myid = out.ID
	getIpnsAdrs()
	return nil
}*/
func IpfsGet(hash, path string) error {
	if err := chkStat(); err != nil {
		return err
	}
	return shell.Get(hash, path)
}

func IpfsBlockGet(hash string) ([]byte, error) {
	if err := chkStat(); err != nil {
		return nil, err
	}
	return shell.BlockGet(hash)
}

func IpfsAddFile(path string) (string, error) {
	if err := chkStat(); err != nil {
		return "", err
	}
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	return IpfsAdd(file)
}

func IpfsAdd(reader io.Reader) (string, error) {
	if err := chkStat(); err != nil {
		return "", err
	}
	hash, err := shell.Add(reader)
	if err != nil {
		fmt.Printf("ERR**:%s\n", err.Error())
		return "", err
	}
	return hash, nil
}

var dIR_IPFS = "/ipfs/"
var hASH_EMPTY_FILE = "QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH"
var hASH_EMPTY_DIR = "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"

const (
	ERR_IpfsCreateIpfsDir_01 ErrCode = "01" // ipns link is created as a file
	ERR_IpfsCreateIpfsDir_02 ErrCode = "02" // dir which is used for boards of alterorg is created as a file
)

func IpfsCreateBoardDir(name string) error {
	if err := chkNameResolved(); err != nil {
		return err
	}
	obj, err := shell.ObjectGet(nsAdrs)
	if err != nil {
		return err
	}
	if len(obj.Links) == 0 && nsAdrs != dIR_IPFS+hASH_EMPTY_DIR { // link to a file
		return makeError(ERR_IpfsCreateIpfsDir_01)
	} else { // link to a directory
		lst, err := shell.List(nsAdrs)
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
			if err = putIpnsAdrs(newhash); err != nil {
				return err
			}
		}
	}
	return nil
}

func getIpnsAdrs() error {
	if err := chkStat(); err != nil {
		return err
	}
	var err error
	if nsAdrs, err = shell.Resolve(myid); err != nil {
		return err
	}
	logIpfs("IPNS:%s", nsAdrs)
	return nil
}

func putIpnsAdrs(adrs string) error {
	if err := chkStat(); err != nil {
		return err
	}
	nsAdrs = adrs
	logIpfs("Start publishing Ipns Address :%s", adrs)
	if err := shell.Publish("", adrs); err != nil {
		return err
	}
	logIpfs("New Ipns Address is :%s", adrs)
	return nil
}

const (
	ERR_IpfsWriteToBoard_01 ErrCode = "01" // dir for boards is not created yet
	ERR_IpfsWriteToBoard_02 ErrCode = "02" // failed to get dir for boards
)

func IpfsWriteToBoard(dir string, data string, n string) error {

	if err := chkStat(); err != nil {
		return err
	}
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
	logIpfs("BoardHash : %s", hash)
	// boarddir
	obj, err := shell.ObjectGet(nsAdrs + "/" + dir)
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
	logIpfs("NewHash for boardforalterorg : %s", hash)

	// ipnsdir
	nsobj, err := shell.ObjectGet(nsAdrs)
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
	if err = putIpnsAdrs(nwnshash); err != nil {
		return err
	}
	return nil
}

func IpfsListBoard(dir string) ([][]string, error) {
	// need to append lists of other nodes
	ret := [][]string{}
	if err := chkStat(); err != nil {
		return nil, err
	}
	list, err := shell.List(nsAdrs + "/" + dir)
	if err != nil {
		return nil, err
	}
	for _, lnk := range list {
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
