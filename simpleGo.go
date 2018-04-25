package main

import (
	"runtime"
	"fmt"
	"time"
)

func say(s string)  {
	for i := 0; i<5 ;i++  {
		runtime.Gosched() // 出让时间片，防止锁死CPU
		println(s)
	}
}
func showLoading(delay int)  {
	var d time.Duration= 100 * time.Millisecond
	for{
		for _, r:=range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(d)
		}
	}
}

func fibc(x int) int {
	if(x < 2){
		return x
	}
	return fibc(x-1) + fibc(x-2)
}

func main()  {
	//go say("hello")
	//go say("pep")
	//say("world")

	go showLoading(100) // 清屏显示loading
	result := fibc(45) // 同步进行实际计算
	fmt.Printf("fibc result 45: %d", result)
}
