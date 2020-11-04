package client

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

var version = "v1.0.0"

// BCacheClient object
type BCacheClient struct {
	addr string
}

// NewCache returns a new cache handler
func NewCache(addr string) *BCacheClient {
	return &BCacheClient{
		addr: "http://" + addr + "/v1/",
	}
}

// Set a new k/v
func (b *BCacheClient) Set(bucket, key, value string) error {
	u := b.addr + bucket + "/" + key

	req, err := http.NewRequest("PUT", u, bytes.NewBuffer([]byte(value)))
	req.Header.Add("User-Agent", "BCache client "+version)

	// Create a new client
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// Get a k/v pair
func (b *BCacheClient) Get(bucket, key string) ([]byte, error) {
	u := b.addr + bucket + "/" + key

	req, err := http.NewRequest("GET", u, nil)
	req.Header.Add("User-Agent", "BCache client "+version)

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

// Del a k/v pair
func (b *BCacheClient) Del(bucket, key string) error {
	u := b.addr + bucket + "/" + key

	req, err := http.NewRequest("DELETE", u, nil)
	req.Header.Add("User-Agent", "BCache client "+version)

	// Create a new client
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// GetBuckets stored
func (b *BCacheClient) GetBuckets() ([]byte, error) {
	u := b.addr + "/buckets"

	req, err := http.NewRequest("GET", u, nil)
	req.Header.Add("User-Agent", "BCache client "+version)

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

// GetKeys stored
func (b *BCacheClient) GetKeys(bucket string) ([]byte, error) {
	u := b.addr + "/keys/" + bucket

	req, err := http.NewRequest("GET", u, nil)
	req.Header.Add("User-Agent", "BCache client "+version)

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
