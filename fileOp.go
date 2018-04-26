package main

import (
	"os"
	"bufio"
	"io"
	"fmt"
	"strings"
)

func main()  {
	f, err := os.Open("./tmp/tinyMps.bak.log")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	startN := 0
	endN := 0
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}
		if strings.Contains(line, "startG"){
			startN ++
		}
		if strings.Contains(line, "endG"){
			endN ++
		}
		liveN := startN - endN
		fmt.Println("current liveN :" , liveN)
		if strings.Contains(line, "current goroutin"){
			fmt.Println(line)
		}
	}
}