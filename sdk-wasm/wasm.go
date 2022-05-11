package main

import (
	"bulwark-vault/sdk"
	"fmt"
	"syscall/js"
)

func main() {
	c := make(chan struct{})
	sdk.NewVault()
	fmt.Println("test")
	js.Global().Set("foo", js.ValueOf("bar"))
	<-c
}
