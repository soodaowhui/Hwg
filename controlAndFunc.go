package main
// 控制语句与函数
import "fmt"

func ifFunc(x int)  { // 类型在参数名后面
	fmt.Println("ifFunc")
	if x>5{
		fmt.Println("x>5")
	}else{
		fmt.Println("x<=5")
	}
}

func muliReturn(x int)(return1 int, return2 bool) { //返回值类型也要指定，且可以指多个
	fmt.Println("muliReturn")
	return x^2, (x>5)
}

func main()  {
	i := 2
	x := 2^i
	ifFunc(x)
	muliReturn(x)
}