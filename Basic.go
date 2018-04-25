package main

import (
	"fmt"
	"strings"
)

func localv()  {
	var v1 int = 1 // var 变量名 变量类型 = 变量值
	v2 := 2 // 上面的简写
	_, v4 := 3, 4 // _ 被丢弃的值
	fmt.Println(v1)
	fmt.Println(v2)
	fmt.Println(v4)
}

func stringOperation()  {
	s := "something here"
	fmt.Println(s + " and more")
	fmt.Println("length of s:" , len(s))
	splitArr := strings.Fields(s)
	fmt.Println(splitArr)
	fmt.Println("length of splitArr: ", len(splitArr))

	for i,v := range(splitArr){
		fmt.Println(i, v)
	}

	muliLineStr := `more
		rows`
	fmt.Println(muliLineStr)
}

func muliv()  {
	var(
		i int
		pi float32
	)
	i = 1
	pi = 3.2
	fmt.Println(float32(i)+pi) // 需要强制类型转换
}

func makev()  { // make关键词用于内建类型的内存分配

}

func main()  {
	localv()
	stringOperation()
	muliv()
	makev()
}