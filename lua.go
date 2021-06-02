package request

import "github.com/rock-go/lua"

func LuaInjectApi(G *lua.UserKV) {
	G.Set("request" , injectHttpRequest())
}