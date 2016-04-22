// cli/rpc.go
// This source file provides function for client of JSONRPC2.0.
//
// Copyright holder is set forth in alterorg.md

package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type jsonreq struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      uint64      `json:"id"`
}

type jsonres struct {
	Version string           `json:"jsonrpc"`
	Result  *json.RawMessage `json:"result"`
	Error   *json.RawMessage `json:"error"`
	Id      uint64           `json:"id"`
}

type Unknown interface {
}

type errorCode int

/*
func newfileUploadRequest(uri string, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body, "-----123445")
	part, err := writer.CreateFormFile("path", path)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	//	for key, val := range params {
	//		_ = writer.WriteField(key, val)
	//	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return http.NewRequest("POST", uri, body)
}

//func AddFile(url string, data io.Reader, reply Unknown) error {
func AddFile(url string, path string, reply Unknown) error {

	//values := neturl.Values{}
	//values.Set("path", path)
	//req, err := http.NewRequest(
	//	"POST",
	//	url,
	//	strings.NewReader(values.Encode()),
	//)
	fmt.Printf("path:%s\n", path)
	req, err := newfileUploadRequest(url, path)
	if err != nil {
		fmt.Printf("falure to create req:%s\n", err.Error())
		return err
	}
	// TODO:implements correct boundary
	req.Header.Set("Content-Type", "multipart/form-data; boundary=-----123445")
	//req.Header.Set("Content-Type", "multipart/form-data;")
	req.Header.Set("Content-Disposition", "form-data: name=\"files\"")

	client := &http.Client{Timeout: 15 * time.Second,
		Transport: &http.Transport{DisableKeepAlives: true},
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	contentType := resp.Header.Get("Content-Type")
	contents := strings.Split(contentType, ";")

	// TODO:implements all type

	switch contents[0] {
	case "application/json":
		out, _ := ioutil.ReadAll(resp.Body)
		return errors.New("testError!!!!:" + string(out))
		var c jsonres
		if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
			fmt.Printf("EEE1:%s\n", err.Error())
			return err
		}
		if c.Error != nil {
			fmt.Printf("EEE2:%s\n", string(*c.Error))
			return errors.New(string(*c.Error))
		}
		if c.Result == nil {
			fmt.Print("EEE2:NIL\n")
			// TODO:!!
			data, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("geti2(%s) -> %s\n", url, string(data))
			return nil
		}
		fmt.Printf("get(%s) -> %s\n", url, c.Result)
		err = json.Unmarshal(*c.Result, reply)
		if err != nil {
			return err
		}
		return nil
	case "text/plain":
		out, _ := ioutil.ReadAll(resp.Body)
		return errors.New("error returned from IPFS:" + string(out))
	default:
		return errors.New("not support response type yet\n" + contentType)
	}
	return nil
}*/

func Request(url string, method string, args []Unknown, reply Unknown) error {
	jsdata, _ := json.Marshal(&jsonreq{
		Version: "2.0",
		Method:  method,
		Params:  args,
		Id:      uint64(rand.Int63()),
	})
	fmt.Print(string(jsdata) + "\n")
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsdata),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var c jsonres
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return err
	}
	if c.Error != nil {
		return errors.New(string(*c.Error))
	}
	if c.Result == nil {
		return nil
	}
	fmt.Printf("request(%s) -> %s\n", method, c.Result)
	err = json.Unmarshal(*c.Result, reply)
	if err != nil {
		return err
	}
	return nil
}
