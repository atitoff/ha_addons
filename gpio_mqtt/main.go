package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println(len(os.Args), os.Args)
	for {
		time.Sleep(10 * time.Second)
		fmt.Println("tick")
	}
}
