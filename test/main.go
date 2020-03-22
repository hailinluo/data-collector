package main

import (
	"fmt"
	"strings"
)

func main() {
	scale := "跟踪标的：中证红利低波动100指数 | 跟踪误差：--"
	strs := strings.Split(scale, "跟踪误差：")
	for _, str := range strs {
		fmt.Printf("[%s]\n", str)
	}

	return
	scale = strings.TrimSuffix(scale, " ")
	fmt.Printf("[%s]\n", scale)
	scale = strings.TrimSuffix(scale, "|")
	fmt.Printf("[%s]\n", scale)
	scale = strings.TrimSuffix(scale, " ")
	fmt.Printf("[%s]\n", scale)
}
