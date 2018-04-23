package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/yanyiwu/gojieba"
)

func TestFirst(t *testing.T) {
	fmt.Println("hello world")
	t.Log("echo hello world")
}

func TestJieba(t *testing.T) {
	var jb = gojieba.NewJieba()
	var str = "我来到北京清华大学"
	var words = jb.CutAll(str)
	fmt.Println(str)
	fmt.Println("全模式:", strings.Join(words, "///"))
}
