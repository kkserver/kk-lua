package lua

/*
#cgo LDFLAGS: -ldl -llua

#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>
#include <stdlib.h>
#include "golua.h"
*/
import "C"

import (
	"github.com/kkserver/kk-lib/kk/dynamic"
	"reflect"
	"unsafe"
)

var g = NewGMutexObject()

type Function func(L *State) int

type Invoke interface {
	Call(L *State) int
}

type State struct {
	id      *C.struct_lua_State
	objects *GObject
}

func NewState() *State {
	v := State{}
	v.id = C.luaL_newstate()
	v.objects = NewGObject()
	g.Set(GObjectId(v.id), &v)
	return &v
}

func (L *State) Close() {
	C.lua_close(L.id)
	g.Remove(GObjectId(L.id))
}

func (L *State) Openlibs() {
	C.luaL_openlibs(L.id)
}

func (L *State) GetTop() int {
	return int(C.lua_gettop(L.id))
}

func (L *State) GetType(idx int) LuaValType {
	return LuaValType(C.lua_type(L.id, C.int(idx)))
}

func (L *State) IsInteger(idx int) bool {
	if L.IsNumber(idx) {
		v := L.ToNumber(idx)
		if float64(int64(v)) == v {
			return true
		}
	}
	return false
}

func (L *State) IsNumber(idx int) bool {
	return C.lua_isnumber(L.id, C.int(idx)) != 0
}

func (L *State) IsBoolean(idx int) bool {
	return C.lua_type(L.id, C.int(idx)) == C.LUA_TBOOLEAN
}

func (L *State) IsString(idx int) bool {
	return C.lua_isstring(L.id, C.int(idx)) != 0
}

func (L *State) IsObject(idx int) bool {
	return C.lua_type(L.id, C.int(idx)) == C.LUA_TUSERDATA
}

func (L *State) IsFunction(idx int) bool {
	return C.lua_type(L.id, C.int(idx)) == C.LUA_TFUNCTION
}

func (L *State) ToObject(idx int) interface{} {
	switch C.lua_type(L.id, C.int(idx)) {
	case C.LUA_TUSERDATA:
		id := C.lua_touserdata(L.id, C.int(idx))
		return L.objects.Get(GObjectId(id))
	case C.LUA_TNUMBER:
		if L.IsInteger(idx) {
			return L.ToInteger(idx)
		}
		return L.ToNumber(idx)
	case C.LUA_TBOOLEAN:
		return L.ToBoolean(idx)
	case C.LUA_TSTRING:
		return L.ToString(idx)
	case C.LUA_TTABLE:
		vs := []interface{}{}
		m := map[interface{}]interface{}{}
		i := int64(0)
		size := int64(0)

		C.lua_pushnil(L.id)

		for L.Next(idx-1) != 0 {

			switch C.lua_type(L.id, -2) {
			case C.LUA_TNUMBER:
				if i+1 == L.ToInteger(-2) {
					i = i + 1
				}
				v := L.ToObject(-1)
				if v != nil {
					vs = append(vs, v)
				}
			case C.LUA_TSTRING:
				key := L.ToString(-2)
				v := L.ToObject(-1)
				if v != nil {
					m[key] = v
				}
			}

			size = size + 1

			L.Pop(1)
		}

		if size == 0 {
			return m
		}

		if i == size {
			return vs
		}

		return m
	}
	return nil
}

func (L *State) ToObjectId(idx int) GObjectId {
	if C.lua_type(L.id, C.int(idx)) == C.LUA_TUSERDATA {
		id := C.lua_touserdata(L.id, C.int(idx))
		return GObjectId(id)
	}
	return nil
}

func (L *State) ToInteger(idx int) int64 {
	return int64(C.lua_tointegerx(L.id, C.int(idx), nil))
}

func (L *State) ToNumber(idx int) float64 {
	return float64(C.lua_tonumberx(L.id, C.int(idx), nil))
}

func (L *State) ToBoolean(idx int) bool {
	return C.lua_toboolean(L.id, C.int(idx)) != 0
}

func (L *State) ToString(idx int) string {
	var size C.size_t = 0
	r := C.lua_tolstring(L.id, C.int(idx), &size)
	return C.GoStringN(r, C.int(size))
}

func (L *State) NewTable() {
	C.lua_createtable(L.id, 0, 0)
}

func (L *State) RawSet(idx int) {
	C.lua_rawset(L.id, C.int(idx))
}

func (L *State) RawGet(idx int) {
	C.lua_rawget(L.id, C.int(idx))
}

func (L *State) RawGeti(idx int, ref int) {
	C.lua_rawgeti(L.id, C.int(idx), C.lua_Integer(ref))
}

func (L *State) Ref(idx int) int {
	return int(C.luaL_ref(L.id, C.int(idx)))
}

func (L *State) UnRef(idx int, ref int) {
	C.luaL_unref(L.id, C.int(idx), C.int(ref))
}

func (L *State) PushString(value string) {
	Cstr := C.CString(value)
	defer C.free(unsafe.Pointer(Cstr))
	C.lua_pushlstring(L.id, Cstr, C.size_t(len(value)))
}

func (L *State) PushInteger(value int64) {
	C.lua_pushinteger(L.id, C.lua_Integer(value))
}

func (L *State) PushNumber(value float64) {
	C.lua_pushnumber(L.id, C.lua_Number(value))
}

func (L *State) PushBoolean(value bool) {
	if value {
		C.lua_pushboolean(L.id, 1)
	} else {
		C.lua_pushboolean(L.id, 0)
	}
}

func (L *State) PushNil(value bool) {
	C.lua_pushnil(L.id)
}

func (L *State) PushFunction(fn Function) {
	id := C.clua_newfunction(L.id)
	L.objects.Set(GObjectId(id), fn)
}

func (L *State) PushObject(object interface{}) {

	if object == nil {
		C.lua_pushnil(L.id)
	} else {
		switch object.(type) {
		case string:
			L.PushString(object.(string))
		case int8:
			L.PushInteger(int64(object.(int8)))
		case int16:
			L.PushInteger(int64(object.(int16)))
		case int32:
			L.PushInteger(int64(object.(int32)))
		case int64:
			L.PushInteger(object.(int64))
		case int:
			L.PushInteger(int64(object.(int)))
		case uint8:
			L.PushInteger(int64(object.(uint8)))
		case uint16:
			L.PushInteger(int64(object.(uint16)))
		case uint32:
			L.PushInteger(int64(object.(uint32)))
		case uint64:
			L.PushInteger(int64(object.(uint64)))
		case uint:
			L.PushInteger(int64(object.(uint)))
		case float32:
			L.PushNumber(float64(object.(float32)))
		case float64:
			L.PushNumber(object.(float64))
		case bool:
			L.PushBoolean(object.(bool))
		default:
			if reflect.TypeOf(object).Kind() == reflect.Func {
				id := C.clua_newfunction(L.id)
				L.objects.Set(GObjectId(id), object)
				return
			}
			{
				_, ok := object.(Invoke)
				if ok {
					id := C.clua_newfunction(L.id)
					L.objects.Set(GObjectId(id), object)
					return
				}
			}
			id := C.clua_newobject(L.id)
			L.objects.Set(GObjectId(id), object)
		}
	}
}

func (L *State) PushValue(idx int) {
	C.lua_pushvalue(L.id, C.int(idx))
}

func (L *State) Next(idx int) int {
	return int(C.lua_next(L.id, C.int(idx)))
}

func (L *State) Pop(n int) {
	C.lua_settop(L.id, C.int(-n-1))
}

func (L *State) SetMetaTable(idx int) int {
	return int(C.lua_setmetatable(L.id, C.int(idx)))
}

func (L *State) SetGlobal(name string) {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))
	C.lua_setglobal(L.id, Cname)
}

func (L *State) LoadString(code string) int {
	Ccode := C.CString(code)
	defer C.free(unsafe.Pointer(Ccode))
	return int(C.luaL_loadstring(L.id, Ccode))
}

func (L *State) LoadFile(path string) int {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))
	return int(C.luaL_loadfilex(L.id, Cpath, nil))
}

func (L *State) Call(n int, r int) int {
	return int(C.lua_pcallk(L.id, C.int(n), C.int(r), C.int(0), C.lua_KContext(0), nil))
}

//export g_object_delete
func g_object_delete(state unsafe.Pointer) {
	v := g.Get(GObjectId(state))
	if v != nil {
		s, ok := v.(*State)
		if ok {
			top := s.GetTop()
			if top > 0 && s.GetType(-top) == LUA_TOBJECT {
				s.objects.Remove(s.ToObjectId(-top))
			}
		}
	}
}

//export g_object_call
func g_object_call(state unsafe.Pointer) int {

	v := g.Get(GObjectId(state))

	if v != nil {
		s, ok := v.(*State)
		if ok {
			idx := int(C.LUA_REGISTRYINDEX - 1)
			if s.GetType(idx) == LUA_TOBJECT {
				v = s.objects.Get(s.ToObjectId(idx))
				if v != nil {
					vv := reflect.ValueOf(v)
					if vv.Kind() == reflect.Func {
						rs := vv.Call([]reflect.Value{reflect.ValueOf(s)})
						if len(rs) > 0 {
							return int(rs[0].Int())
						}
						return 0
					} else {
						invoke, ok := v.(Invoke)
						if ok {
							return invoke.Call(s)
						}
					}
				}
			}

		}
	}

	return 0
}

//export g_object_get
func g_object_get(state unsafe.Pointer) int {

	v := g.Get(GObjectId(state))

	if v != nil {
		s, ok := v.(*State)
		if ok {
			top := s.GetTop()
			if top > 1 && s.GetType(-top) == LUA_TOBJECT && s.IsString(-top+1) {
				v = s.ToObject(-top)
				if v != nil {
					key := s.ToString(-top + 1)
					v = dynamic.Get(v, key)
					s.PushObject(v)

					return 1
				}
			}
		}
	}

	return 0
}

//export g_object_set
func g_object_set(state unsafe.Pointer) int {

	v := g.Get(GObjectId(state))

	if v != nil {
		s, ok := v.(*State)
		if ok {
			top := s.GetTop()
			if top > 1 && s.GetType(-top) == LUA_TOBJECT && s.IsString(-top+1) {
				v = s.ToObject(-top)
				if v != nil {
					key := s.ToString(-top + 1)
					if top > 2 {
						dynamic.Set(v, key, s.ToObject(-top+2))
					} else {
						dynamic.Set(v, key, nil)
					}

				}
			}
		}
	}

	return 0
}
