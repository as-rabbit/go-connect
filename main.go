package main

import (
	"fmt"
	"go-connect/connecter/mysql/file"
)

func main() {
	
	fileConf := file.NewConfig()
	
	d, err := fileConf.Get("test")
	
	fmt.Println(d, err)
	
}
