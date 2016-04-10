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
