package sdk

import (
	"fmt"
	"syscall/js"
)


func await(awaitable js.Value) (*js.Value, *js.Value) {
    then := make(chan *js.Value)
    defer close(then)
    thenFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		arg := args[0]
        then <- &arg
        return nil
    })
    defer thenFunc.Release()

    catch := make(chan *js.Value)
    defer close(catch)
    catchFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Caught error in await")
        catch <- &args[0]
        return nil
    })
    defer catchFunc.Release()

	awaitable.Call("then", thenFunc).Call("catch",catchFunc)


    select {
    case result := <-then:
        return result, nil
    case err := <-catch:
        return nil, err
    }
}