package main

import (
	"log"
	"os"
	"strconv"
)

// Env is wrapper of all sorts db clients
type Env struct {
	S Storage
}

func getEnv() *Env {
	addr := os.Getenv("APP_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	passwd := os.Getenv("APP_REDIS_PASSWD")
	if passwd == "" {
		passwd = ""
	}
	dbS := os.Getenv("APP_REDIS_DB")
	if dbS == "" {
		dbS = "0"
	}
	db, err := strconv.Atoi(dbS)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connect to redis (addr: %s password: %s db: %d)", addr, passwd, db)

	r := NewRedisCli(addr, passwd, db)
	return &Env{S: r}
}
