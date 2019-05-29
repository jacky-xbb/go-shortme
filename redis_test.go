package main_test

import (
	"testing"

	"."
)

var (
	r  *main.RedisCli
	sl string
)

const (
	url = "https://wwww.example.com"
)

func TestNewRedisCli(t *testing.T) {
	addr := "127.0.0.1:6379"
	pass := ""
	db := 0

	r = main.NewRedisCli(addr, pass, db)
	if err := r.Ping(); err != nil {
		t.Errorf("No server running at localhost:6379")
	}
}

func TestShorten(t *testing.T) {
	var err error
	exp := 60
	sl, err = r.Shorten(url, int64(exp))
	if err != nil {
		t.Error("Shorten failed")
	}
}

func TestUnshorten(t *testing.T) {
	u, err := r.Unshorten(sl)
	if err != nil {
		t.Error("Unshorten return error")
	}

	if u != url {
		t.Errorf("Expected receive url %s. Got url %s", url, u)
	}
}
