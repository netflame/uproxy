package main

import (
	"testing"
)

var cli *client

func init() {
	rp := getRedisPool()
	cli = &client{conn: rp.Get()}
}

func TestSet(t *testing.T) {
	if err := cli.set("redigo", "redis"); err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	TestSet(t)
	v, err := cli.get("redigo")
	if err != nil {
		t.Error(err)
	}
	if v != "redis" {
		t.Error("Test Failed")
	}
}

func TestExists(t *testing.T) {
	TestSet(t)
	y, err := cli.exists("redigo")
	if err != nil {
		t.Error(err)
	}
	if !y {
		t.Error("Test Failed")
	}
}

func TestDel(t *testing.T) {
	TestSet(t)
	err := cli.del("redigo")
	if err != nil {
		t.Error(err)
	}
	y, _ := cli.exists("redigo")
	if y {
		t.Error("Test Failed")
	}
}
