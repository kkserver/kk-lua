#include "lua.h"
#include "lauxlib.h"
#include "lualib.h"
#include "golua.h"
#include "_cgo_export.h"

static int clua_object_gc(struct lua_State *L) {
	g_object_delete((uintptr) L);
	return 0;
}

static int clua_object_get(struct lua_State *L) {
	return g_object_get((uintptr) L);
}

static int clua_object_set(struct lua_State *L) {
	return g_object_set((uintptr) L);
}

static int clua_object_call(struct lua_State *L) {
	return g_object_call((uintptr) L);
}

uintptr clua_newobject(struct lua_State *L) {
	
	void * v = lua_newuserdata(L,sizeof(uintptr));

	lua_newtable(L);

	lua_pushstring(L,"__gc");
	lua_pushcfunction(L,clua_object_gc);
	lua_rawset(L,-3);

	lua_pushstring(L,"__index");
	lua_pushcfunction(L,clua_object_get);
	lua_rawset(L,-3);

	lua_pushstring(L,"__newindex");
	lua_pushcfunction(L,clua_object_set);
	lua_rawset(L,-3);

	lua_setmetatable(L,-2);

	return (uintptr) v;
}

uintptr clua_newfunction(struct lua_State *L) {
	
	void * v = lua_newuserdata(L,sizeof(uintptr));

	lua_newtable(L);

	lua_pushstring(L,"__gc");
	lua_pushcfunction(L,clua_object_gc);
	lua_rawset(L,-3);

	lua_pushstring(L,"__index");
	lua_pushcfunction(L,clua_object_get);
	lua_rawset(L,-3);

	lua_pushstring(L,"__newindex");
	lua_pushcfunction(L,clua_object_set);
	lua_rawset(L,-3);

	lua_setmetatable(L,-2);

	lua_pushcclosure(L, clua_object_call, 1);

	return (uintptr) v;
}



