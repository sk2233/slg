/*
@author: sk
@date: 2022/12/25
*/
package test

import (
	"fmt"
	"testing"
)

func Test11(t *testing.T) {
	temp := GetRes()
	if temp == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(temp)
	}
}

func GetRes() []string {
	return []string{}
}

func Test12(t *testing.T) {
	//arr := []string{"sdsds"}
	//arr=append(arr,arr)
	type Name interface {
		Name() string
	}
	arr := make([]string, 0)
	//arr=append(arr,arr)
	iarr := make([]Name, 0)
	//iarr = append(iarr, iarr)
	fmt.Println(arr, iarr)
}

func TestM(t *testing.T) {
	arr := make([]string, 0)
	arr = append(arr, "sdfsdf")
	removeIt(arr)
	fmt.Println(arr)
}

func removeIt(arr []string) {
	arr = make([]string, 0)
}
