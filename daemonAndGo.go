package main

// 微型MPS程序，目标：
// 主程序=原主daemon，控制最大启动协程数（如200个）
// 每个协程只做pop data/转发这两件事，做完就结束
// 有代码更新时，restart主程序，则主程序平滑重启，用新binary启动新进程，老进程仍然等待处理完剩下的内容为止再结束
// 外部命令可控制deamon的启，停，重启等

import (
	"time"
"math/rand"
	"log"
	"os"
	"runtime"
	"bytes"
	"strconv"
	"fmt"
)

const maxRouNum  int = 20 // 最大启动协程数

var debugLog *log.Logger

func fakeRedisRead()  {// 模拟一个redis读取所用的时间
	rand.Seed(time.Now().UnixNano())
	randInt := rand.Intn(20)
	time.Sleep( time.Duration(randInt) * time.Millisecond)
}

func fakeRedisWrite()  { // 模拟一个redis写入所用的时间
	rand.Seed(time.Now().UnixNano())
	randInt := rand.Intn(30)
	time.Sleep( time.Duration(randInt) * time.Millisecond)
}

func doQuery()  { // 处理数据并写入下级队列
	debugLog.Println("goroutine "+GetGID()+": doQuery" )
	fakeRedisWrite()
}

func getData()  { // 获取数据
	debugLog.Println("goroutine "+GetGID()+": getData" )
	fakeRedisRead()
}

func startG(pool chan struct{})  { // 启动一个处理协程
	pool <- struct{}{}// 锁定池中一个资源
	debugLog.Println("goroutine "+GetGID()+": startG" )
	getData()
	doQuery()
	debugLog.Println("goroutine "+GetGID()+": endG" )
	<-pool // 释放池中一个资源
}

func startDaemon()  { // 开启主进程
	pool := make(chan struct{}, maxRouNum) // 锁定池，保证最大启动协程数
	run := true // 是否继续运行，从外部接收命令
	for i := 0; i<1000 ; i++  { //最大循环1000次
		go startG(pool)
		gNum := strconv.Itoa((runtime.NumGoroutine()))
		debugLog.Println("current goroutin number:", gNum, ",length of pool:", strconv.Itoa(len(pool)))
		fmt.Println("current goroutin number:", gNum, ",length of pool:", strconv.Itoa(len(pool)))
		time.Sleep( time.Duration(1) * time.Millisecond) // 如果没有sleep，主程序结束太快了，则所有协程也都一起结束了
		if !run {
			break
		}
	}
}

func GetGID() string { // 获取当前协程号
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return strconv.Itoa(int(n))
}

func main()  {
	fileName := "tmp/tinyMps.log"
	logFile,err  := os.OpenFile(fileName,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error")
	}
	debugLog = log.New(logFile,"[Info]",log.Ldate | log.Ltime)
	//debugLog.Println("A Info message here22")
	//debugLog.SetPrefix("[Debug]")
	//debugLog.Println("A Debug Message here33")
	debugLog.Println("main process start")


	startDaemon()
}
