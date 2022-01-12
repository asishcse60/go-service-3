package main

import (
	"fmt"

	"github.com/asishcse60/service/app/tooling/sales-admin/commands"
)

func main() {
	if err:=commands.GenKey(); err!=nil{
		fmt.Println(err)
	}
}