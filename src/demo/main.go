package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file,err:=os.OpenFile("./../.gitignore",os.O_APPEND,0666)
	if err!=nil {
		fmt.Println(err)
	}
	defer file.Close()
	r:=bufio.NewWriter(file)
	n,err:=r.WriteString("\nyes")
	r.Flush()
	fmt.Println(n,err)
}
