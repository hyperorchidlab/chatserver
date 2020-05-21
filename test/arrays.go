package main

import (
	"github.com/kprc/chatserver/ed25519"
	"fmt"
)

func main()  {
	a:=[]byte{'a'}
	b:=[]byte{'b'}
	c:=[]byte{'c'}

	arr:=make([][]byte,0)
	arr = append(arr,c)
	arr = append(arr,b)
	arr = append(arr,a)

	for i:=0;i<len(arr);i++{
		fmt.Println(string(arr[i]))
	}


	r := chatcrypt.InsertionSortDArray(arr)

	for i:=0;i<len(r);i++{
		fmt.Println(string(r[i]))
	}

}



