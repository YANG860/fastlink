package main

import (
	"log"
	"time"
)

func statQps() {
	var lastCount uint64 = 0
	for {
		time.Sleep(time.Second * 15)
		log.Println("QPS: ", float32(totalQuery-lastCount)/15)
		lastCount = totalQuery
	}
}
