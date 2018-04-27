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

func worker()  { // 启动一个处理协程
	debugLog.Println("goroutine "+GetGID()+": start worker" )
	getData()
	doQuery()
	debugLog.Println("goroutine "+GetGID()+": end worker" )
}


func worker2(ch chan string)  { // 启动一个处理协程
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

func wait(wg sync.WaitGroup, pool chan struct{}, abort chan struct{}) (endLoopNum int,stopEnd bool){
	loopNum := 1000
	for i := 0; i<loopNum ; i++  {
		tick := time.Tick(5 * time.Millisecond)
		wg.Add(1)
		select {
		case <-tick: // 等待1ms
			if i%10 == 0 {LogL("tick")}
			if i%50 == 0 {
				printCurrentNumGo(len(pool))
			}
			go func() {
				pool <- struct{}{}// 锁定池中一个资源
				defer func() {
					<-pool // 确保释放池中一个资源
				}()
				//Worker() // 带超时方案
				worker() // 无超时方案
				wg.Done()
			}()

		case <- abort:
			LogL("get abort signle")
			return i,true //中断返回
		}
	}
	return loopNum, false // 循环结束返回
}

func startDaemon()  { // 开启主进程
	pool := make(chan struct{}, maxRouNum) // 锁定池，限制最大启动协程数

	abort := make(chan struct{})
	sig := -1000
	go func() {
		sig,_ = os.Stdin.Read(make([]byte, 1)) // 从屏幕输入读取一个信号
		abort <- struct{}{}
	}()


	var wg sync.WaitGroup
	endLoopNum, stop :=wait(wg, pool, abort)

	if stop {
		LogL("stop loop at loop num:" + strconv.Itoa(endLoopNum))
		LogL("abort signle:" + strconv.Itoa(sig))
	} else{
		LogL("loop end")
	}
	LogL("current goroutin number:"+ strconv.Itoa((runtime.NumGoroutine()))+ ",length of pool:"+ strconv.Itoa(len(pool)))
	LogL("wait for all goroutin end")
	wg.Wait() // 等待所有启动但还未完成的协程执行完毕
	LogL("current goroutin number:"+ strconv.Itoa((runtime.NumGoroutine()))+ ",length of pool:"+ strconv.Itoa(len(pool)))
	LogL("Daemon end")
}

func printCurrentNumGo(poolNum int)  {
	gNum := strconv.Itoa((runtime.NumGoroutine()))
	LogL("current goroutin number:"+ gNum+ ",length of pool:"+ strconv.Itoa(poolNum))
}

func LogL(s string)  {
	debugLog.Println(s)
	fmt.Println(s)
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
