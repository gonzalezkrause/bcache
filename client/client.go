package data

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

var version = "v0.0.1"

type BoltCacheClient struct {
	addr string
}

func NewCache(addr string) *BoltCacheClient {
	return &BoltCacheClient{
		addr: "http://" + addr + "/v1/",
	}
}

func (b *BoltCacheClient) Set(bucket, key, value string) error {
	u := b.addr + bucket + "/" + key

	fmt.Println(u)

	req, err := http.NewRequest("POST", u, bytes.NewBuffer([]byte(value)))
	req.Header.Add("User-Agent", "BoltCache client "+version)

	// Create a new client
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (b *BoltCacheClient) Get(bucket, key string) ([]byte, error) {
	u := b.addr + bucket + "/" + key

	req, err := http.NewRequest("GET", u, nil)
	req.Header.Add("User-Agent", "BoltCache client "+version)

	// Create a new client
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	// Readout the body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func (b *BoltCacheClient) Del(bucket, key string) error {
	u := b.addr + bucket + "/" + key

	req, err := http.NewRequest("DELETE", u, nil)
	req.Header.Add("User-Agent", "BoltCache client "+version)

	// Create a new client
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
