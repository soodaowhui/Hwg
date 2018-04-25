package main

import "fmt"

func stringAndArr(){
	str1 := "oneString"
	strArr := []byte(str1)
	fmt.Println("before edit:",str1)
	fmt.Println("strArr:",strArr)
	strArr[0] = 'F'
	str2 := string(strArr)
	fmt.Println("after edit:",str2)
}

func main()  {
	stringAndArr()
}