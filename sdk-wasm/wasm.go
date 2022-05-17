package main

import (
	"bulwark-vault/sdk"
	"fmt"
	"syscall/js"
)

func loginToVault(email, password string) *sdk.VaultAccess {
	vault := sdk.NewVault()
	err := vault.Login(email, password)
	if err != nil {
		return nil
	}
	return vault
}

func get(vault sdk.VaultAccess, path string) string {

}

func put(vault sdk.VaultAccess, path string, data string) {

}

func main() {
	c := make(chan struct{})
	sdk.NewVault()
	fmt.Println("test")
	js.Global().Set("foo", js.ValueOf("bar"))
	<-c
}
