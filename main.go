package main

import (
	"context"
	"fmt"
	"go-connect/connecter/file/mysql"
)

func main() {

	ctx := context.TODO()

	newDb := mysql.NewConnector()

	db, err := newDb.Make(ctx, "test")

	fmt.Println(db, err)

	//mysql.NewConnector()

	//fileConf := file.NewConfig()
	//
	//d, err := fileConf.Get("test")
	//
	//fmt.Println(d, err)

}
