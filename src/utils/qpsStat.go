package utils

import (
	"log"
	"time"
)

var TotalQuery uint64 = 0

// statQps 定时统计并打印 QPS
func StatQps() {
	var lastCount uint64 = 0
	for {
		time.Sleep(time.Second * 15)
		log.Println("QPS: ", float32(TotalQuery-lastCount)/15)
		lastCount = TotalQuery
	}
}
