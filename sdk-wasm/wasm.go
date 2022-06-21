package main

import (
	"bulwark-vault/sdk"
	"fmt"
	"syscall/js"
)

func makeAsync(executor func(js.Value, []js.Value) any) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		executor := js.FuncOf(func(_ js.Value, executorArgs []js.Value) interface{} {
			resolve := executorArgs[0]
			go func() {
				val := executor(this, args)
				resolve.Invoke(val)
			}()
			return nil
		})
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(executor)
	})
}

var vault *sdk.VaultAccess = nil

func login(this js.Value, args []js.Value) any {
	if len(args) != 2 {
		return ""
	}
	email := args[0].String()
	password := args[1].String()
	vault = sdk.NewVault()
	err := vault.Login(email, password)
	if err != nil {
		return fmt.Sprintf("Error logging into vault: %s", err)
	}
	return ""
}

func getMasterSecret(this js.Value, args []js.Value) any {
	if vault == nil {
		return ""
	}
	return vault.MasterSecret()
}

func getKeyDirectory(this js.Value, args []js.Value) any {
	if vault == nil {
		return js.ValueOf(nil)
	}
	jsonData, err := vault.KeyDirectory().Json()
	if err != nil {
		return js.ValueOf(nil)
	}
	return js.ValueOf(jsonData)
}

func put(this js.Value, args []js.Value) any {
	if vault == nil {
		return js.ValueOf(nil)
	}
	path := args[0].String()
	data := args[1].String()
	err := vault.Put(path, data)
	if err != nil {
		return js.ValueOf(nil)
	}
	return js.ValueOf("Success")
}

func get(this js.Value, args []js.Value) any {
	if vault == nil {
		return js.ValueOf(nil)
	}
	path := args[0].String()
	data, err := vault.Get(path)
	if err != nil {
		return js.ValueOf(nil)
	}
	return js.ValueOf(data)
}

func getAuthData(this js.Value, args []js.Value) any {
	publicKey := args[0].String()
	encryptionKey := args[1].String()
	data, err := sdk.GetAuthData(publicKey, encryptionKey)
	if err != nil {
		return js.ValueOf(nil)
	}
	return js.ValueOf(data)
}

func main() {
	c := make(chan struct{})
	fmt.Println("WASM started")
	vaultInterface := make(map[string]interface{})
	vaultInterface["login"] = makeAsync(login)
	vaultInterface["getMasterSecret"] = makeAsync(getMasterSecret)
	vaultInterface["getKeyDirectory"] = makeAsync(getKeyDirectory)
	vaultInterface["get"] = makeAsync(get)
	vaultInterface["put"] = makeAsync(put)
	vaultInterface["getAuthData"] = makeAsync(getAuthData)
	js.Global().Set("vaultInterface", js.ValueOf(vaultInterface))
	<-c
}
