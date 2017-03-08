package main

import (
	"github.com/kkserver/kk-lua/lua"
	"log"
)

func main() {

	L := lua.NewState()

	L.Openlibs()

	L.PushObject(map[interface{}]interface{}{"title": "OK", "onload": func(L *lua.State) int {

		top := L.GetTop()

		for i := 0; i < top; i++ {
			log.Println(L.ToObject(-top + i))
		}

		return 0
	}})

	L.SetGlobal("data")

	if L.LoadFile("./main.lua") == 0 {

		if L.Call(0, 0) != 0 {
			log.Println(L.ToString(-1))
			L.Pop(1)
		}

	} else {
		log.Println(L.ToString(-1))
		L.Pop(1)
	}

	L.Close()

}
