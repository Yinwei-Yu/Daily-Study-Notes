package main
// 一个接口底层占两个字节
// 一个字节存放类型信息，另一个字节存放数据或者指向数据的指针
import (
	"fmt"
)

type Shape interface {
	Area() float32
}

type Square struct {
	side float32
}

type Rectangle struct {
	width float32
	length float32
}

func (s Square)Area() float32 {
	return s.side*s.side
}

func (r Rectangle)Area() float32 {
	return r.length*r.width
}

func main() {
	r:=Rectangle{3,2}
	q:=Square{5.0}
	shapes := []Shape{r,q}
	
	fmt.Println("looping shapes")
	
	for n,_:=range shapes {
		fmt.Printf("Shape type :%T\n", shapes[n])
		fmt.Println("Shape area:",shapes[n].Area())
	}
	
}