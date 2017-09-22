#include <stdint.h>

#ifndef KK_GO_LUA_H
#define KK_GO_LUA_H

typedef void * kk_uintptr;

kk_uintptr clua_newobject(struct lua_State * L);

kk_uintptr clua_newfunction(struct lua_State *L);

#endif
