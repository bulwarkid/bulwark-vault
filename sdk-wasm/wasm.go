package main

import (
	"bulwark-vault/sdk"
	"fmt"
	"syscall/js"
)

func makeAsync(executor func(js.Value, []js.Value) any) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		handler := js.FuncOf(func(_ js.Value, args2 []js.Value) interface{} {
			resolve := args2[0]
			go func() {
				val := executor(this, args)
				resolve.Invoke(val)
			}()
			return nil
		})
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	});
}

func loginToVault(email, password string) (*sdk.VaultAccess, error) {
	vault := sdk.NewVault()
	err := vault.Login(email, password)
	if err != nil {
		return nil, err
	}
	return vault, nil
}

var vault *sdk.VaultAccess = nil

func login(this js.Value, args []js.Value) any {
	if len(args) != 2 {
		return "";
	}
	email := args[0].String()
	password := args[1].String()
	vault = sdk.NewVault()
	err := vault.Login(email, password)
	if err != nil {
		return fmt.Sprintf("Error logging into vault: %s", err)
	}
	return "";
}

func getMasterSecret(this js.Value, args []js.Value) any {
	if vault == nil {
		return ""
	}
	return vault.MasterSecret()
}

func main() {
	c := make(chan struct{})
	fmt.Println("WASM started")
	js.Global().Set("vaultLogin", makeAsync(login))
	js.Global().Set("vaultGetMasterSecret", js.FuncOf(getMasterSecret))
	js.Global().Set("test",js.ValueOf("foo"))
	<-c
}
