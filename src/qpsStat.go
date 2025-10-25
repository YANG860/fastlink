package main

import (
	"log"
	"time"
)

// statQps 定时统计并打印 QPS
func statQps() {
	var lastCount uint64 = 0
	for {
		time.Sleep(time.Second * 15)
		log.Println("QPS: ", float32(totalQuery-lastCount)/15)
		lastCount = totalQuery
	}
}
