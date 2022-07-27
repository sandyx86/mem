package main

import (
	"fmt"
	"os"
)

func main() {
	var a int = 35;
	fmt.Println("the value of a is:", a, &a, os.Getpid())
	fmt.Scanf("Press any key to continue")
}