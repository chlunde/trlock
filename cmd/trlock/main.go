// Package main is a simple demo/test for trlock.  Stop with ctrl-c
package main

import (
	"log"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/chlunde/trlock"
)

func main() {
	X, err := xgb.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	defer X.Close()

	for {
		err = trlock.Lock(X)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Locked")
		time.Sleep(5 * time.Second)

		log.Println("Unlocking")
		trlock.Unlock(X)

		log.Println("Unlocked; ctrl-c now")
		time.Sleep(3 * time.Second)
	}
}
