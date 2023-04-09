package main

import (
	"context"
	"fmt"
	"github.com/as-rabbit/go-connect/connecter/file/mysql"
)

func main() {

	ctx := context.TODO()

	c, err := mysql.NewConfig()

	newDb := mysql.NewConnector(c)

	db, err := newDb.Make(ctx, "test")

	fmt.Println(db, err)

	var count int64

	db.Table("goods").Where("id=?", 1).Count(&count)

	fmt.Println(count)

	//mysql.NewConnector()

	//fileConf := file.NewConfig()
	//
	//d, err := fileConf.Get("test")
	//
	//fmt.Println(d, err)

}
