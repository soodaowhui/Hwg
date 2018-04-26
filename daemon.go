package main

// 重写MPS程序，目标：
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
	"os/exec"
	"sync"
)

const maxRouNum  int = 20 // 最大启动协程数
const workerTimeOut int = 1 // 单worker 1秒超时

var debugLog *log.Logger

func fakeRedisRead()  {// 模拟一个redis读取所用的时间
	rand.Seed(time.Now().UnixNano())
	randInt := rand.Intn(20) // 为体现协程调度的效果，将redis操作时间放大为最大20ms，实际应该在1ms以下
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

func worker(pool chan struct{})  { // 启动一个处理协程
	pool <- struct{}{}// 锁定池中一个资源
	defer func() {
		debugLog.Println("goroutine "+GetGID()+": release pool" )
		<-pool // 确保释放池中一个资源
	}()
	debugLog.Println("goroutine "+GetGID()+": start worker" )
	getData()
	doQuery()
	debugLog.Println("goroutine "+GetGID()+": end worker" )
}


func worker2(pool chan struct{}, ch chan string)  { // 启动一个处理协程
	pool <- struct{}{}// 锁定池中一个资源
	defer func() {
		debugLog.Println("goroutine "+GetGID()+": release pool" )
		<-pool // 确保释放池中一个资源
	}()
	debugLog.Println("goroutine "+GetGID()+": start worker2" )
	getData()
	doQuery()
	debugLog.Println("goroutine "+GetGID()+": end worker2" )
	ch <- "success"
}

//func Worker(pool chan struct{}){ // 包装worker2函数，用于设置超时
//	ch_run := make(chan string)
//	go worker2(pool, ch_run)
//	select {
//	case re := <-ch_run:
//		_ <- re  // 成功结束，不处理
//	case <-time.After(time.Duration(workerTimeOut) * time.Second):
//		return // 超时，结束当前go，同时结束子Go
//	}
//}

func startDaemon()  { // 开启主进程
	pool := make(chan struct{}, maxRouNum) // 锁定池，保证最大启动协程数
	run := true // 是否继续运行，从外部接收命令
	var wg sync.WaitGroup
	for i := 0; i<100 ; i++  { //测试时最大循环100次
		wg.Add(1)
		go func() {
			//Worker(pool) // 带超时方案
			worker(pool) // 无超时方案
			wg.Done()
		}()

		gNum := strconv.Itoa((runtime.NumGoroutine()))
		debugLog.Println("current goroutin number:", gNum, ",length of pool:", strconv.Itoa(len(pool)))
		fmt.Println("current goroutin number:", gNum, ",length of pool:", strconv.Itoa(len(pool)))
		time.Sleep( time.Duration(1) * time.Millisecond) // 如果没有sleep，主程序结束太快了，然后所有协程也都一起结束了
		if !run {
			break
		}
	}
	debugLog.Println("loop end")
	debugLog.Println("current goroutin number:", strconv.Itoa((runtime.NumGoroutine())), ",length of pool:", strconv.Itoa(len(pool)))
	wg.Wait() // 等待所有启动但还未完成的协程执行完毕
	debugLog.Println("current goroutin number:", strconv.Itoa((runtime.NumGoroutine())), ",length of pool:", strconv.Itoa(len(pool)))
	debugLog.Println("Daemon end")
}

// In order to keep the working directory the same as when we started we record
// it at startup.
var originalWD, _ = os.Getwd()

// 重启daemon,载入新代码
func reloadDaemon()  (int, error) {
	// Use the original binary location. This works with symlinks such that if
	// the file it points to has been changed we will use the updated symlink.
	argv0, err := exec.LookPath(os.Args[0])
	if err != nil {
		return 0, err
	}

	// 复制 environment
	var env []string
	for _, v := range os.Environ() {
		env = append(env, v)
	}

	process, err := os.StartProcess(argv0, os.Args, &os.ProcAttr{
		Dir:   originalWD,
		Env:   env,
	})
	if err != nil {
		return 0, err
	}
	return process.Pid, nil
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
	pid := os.Getegid()
	debugLog.Println("main process start, pid:", pid)

	startDaemon()
}
