package main

import (
	"github.com/kprc/chatserver/db"
	"fmt"
)

func main() {
	//a := []byte{'a'}
	//b := []byte{'b'}
	//c := []byte{'c'}
	//
	//arr := make([][]byte, 0)
	//arr = append(arr, c)
	//arr = append(arr, b)
	//arr = append(arr, a)
	//
	//for i := 0; i < len(arr); i++ {
	//	fmt.Println(string(arr[i]))
	//}
	//
	//r := chatcrypt.InsertionSortDArray(arr)
	//
	//for i := 0; i < len(r); i++ {
	//	fmt.Println(string(r[i]))
	//}

	sesc:=db.Discrete2Section([]int{1,4,6,8,9,13,14,15})

	for i:=0;i<len(sesc);i++{
		fmt.Println(sesc[i].String())
	}

}
