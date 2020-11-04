package client_test

import (
	"strconv"
	"testing"

	"github.com/gonzalezkrause/bcache/client"
)

func TestNewCache(t *testing.T) {
	cache := client.NewCache("127.0.0.1:3018")
	if cache == nil {
		t.Fail()
	}
}

func TestSet(t *testing.T) {
	data := "{\"a\":42, \"b\":\"foo\"}"

	cache := client.NewCache("127.0.0.1:3018")

	for i := 0; i < 5; i += 1 {
		for j := 0; j < 5; j += 1 {
			if err := cache.Set("c"+strconv.Itoa(i), "k"+strconv.Itoa(j), data); err != nil {
				t.Error(err)
			}
		}
	}
}

func TestGet(t *testing.T) {
	data := "{\"a\":42, \"b\":\"foo\"}"

	d, err := client.NewCache("127.0.0.1:3018").Get("c1", "k1")
	if err != nil {
		t.Error(err)
	}

	if string(d) != data {
		t.Errorf("Expected: %#v Got: %#v", data, string(d))
	}
}

func TestGetBuckets(t *testing.T) {
	expected := "c0\nc1\nc2\nc3\nc4\n"

	got, err := client.NewCache("127.0.0.1:3018").GetBuckets()
	if err != nil {
		t.Error(err)
	}

	if expected != string(got) {
		t.Errorf("Expected: %#v Got: %#v", expected, string(got))
	}
}

func TestGetKeys(t *testing.T) {
	expected := "k0\nk1\nk2\nk3\nk4\n"

	got, err := client.NewCache("127.0.0.1:3018").GetKeys("c1")
	if err != nil {
		t.Error(err)
	}

	if expected != string(got) {
		t.Errorf("Expected: %#v Got: %#v", expected, string(got))
	}
}

func TestDel(t *testing.T) {
	cache := client.NewCache("127.0.0.1:3018")

	for i := 0; i < 5; i += 1 {
		for j := 0; j < 5; j += 1 {
			if err := cache.Del("c"+strconv.Itoa(i), "k"+strconv.Itoa(j)); err != nil {
				t.Error(err)
			}
		}
	}
}
