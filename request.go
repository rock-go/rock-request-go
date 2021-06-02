package request

import (
	"github.com/rock-go/lua"
	"github.com/go-resty/resty/v2"
)

const (
	MethodGet     string = "GET"
	MethodPatch   string = "PATCH"
	MethodTrace   string = "TRACE"
	MethodOptions string = "OPTIONS"
	MethodDelete  string = "DELETE"
	MethodHead    string = "HEAD"
	MethodPut     string = "PUT"
	MethodPost    string = "POST"
)

type httpResponse struct {
	lua.Super
	rc *resty.Response
	err error
}

func(hr *httpResponse) Name() string {
	return "rock.http.response"
}

func(hr *httpResponse) ToLightUserData() *lua.LightUserData {
	return lua.NewLightUserData(hr)
}

func (hr *httpResponse) Index( L *lua.LState , key string ) lua.LValue {
	if key == "code" { return lua.LNumber(hr.rc.StatusCode()) }
	if key == "body" { return lua.LString(hr.rc.Body())       }
	if key == "err"  { return hr.Error()     }

	return nil
}

func (hr *httpResponse) Error() lua.LValue {
	if hr.err == nil {
		return lua.LNil
	}
	return lua.LString(hr.err.Error())
}


type httpRequest struct {
	lua.Super

	client *resty.Client
	r  *resty.Request
}

func(r *httpRequest) Name() string {
	return "rock.http.request"
}

func (r *httpRequest) ToLightUserData() *lua.LightUserData {
	L := lua.State()
	ud := L.NewLightUserData(r)
	lua.FreeState(L)
	return ud
}

func (r *httpRequest) output(L *lua.LState) int {
	r.r.SetOutput(L.CheckString( 1))
	L.Push(L.NewLightUserData(r))
	return 1
}

func (r *httpRequest) GET( L *lua.LState ) int  {
	return r.Execute(MethodGet, L)
}

func (r *httpRequest) POST( L *lua.LState ) int {
	return r.Execute(MethodPost, L)
}

func (r *httpRequest) PUT( L *lua.LState ) int {
	return r.Execute(MethodPut,L)
}

func (r *httpRequest) OPTIONS( L *lua.LState ) int {
	return r.Execute(MethodOptions,L )
}
func (r *httpRequest) TRACE( L *lua.LState ) int {
	return r.Execute(MethodTrace,L)
}

func (r *httpRequest) HEAD( L *lua.LState ) int {
	return	r.Execute(MethodHead,L)
}

func (r *httpRequest) DELETE( L *lua.LState ) int {
	return	r.Execute(MethodDelete,L)
}

func (r *httpRequest) PATCH( L *lua.LState ) int {
	return	r.Execute(MethodPatch,L)
}
func (r *httpRequest) Execute(method string , L *lua.LState) int {
	rc , err := r.r.Execute( method , L.CheckString(1))
	resp := &httpResponse{ rc: rc , err: err}
	L.Push(resp.ToLightUserData())
	return 1
}

func (r *httpRequest) Index( L *lua.LState , key string ) lua.LValue{ //lua代码获取对象的方法

	if key == "OPTIONS"     { return lua.NewFunction( r.OPTIONS )  }
	if key == "DELETE"      { return lua.NewFunction( r.DELETE  )  }
	if key == "PATCH"       { return lua.NewFunction( r.PATCH   )  }
	if key == "TRACE"       { return lua.NewFunction( r.TRACE   )  }
	if key == "POST"        { return lua.NewFunction( r.POST    )  }
	if key == "HEAD"        { return lua.NewFunction( r.HEAD    )  }
	if key == "GET"         { return lua.NewFunction( r.GET     )  }
	if key == "PUT"         { return lua.NewFunction( r.PUT     )  }

	if key == "output"      { return lua.NewFunction( r.output  )  }

	return lua.LNil
}

func injectHttpRequest() *lua.LightUserData {
	client := resty.New()

	r := &httpRequest{
		client: client,
		r: client.R(),
	}

	return lua.NewLightUserData( r )
}