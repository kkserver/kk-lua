#include <stdint.h>

typedef void * kk_uintptr;

kk_uintptr clua_newobject(struct lua_State * L);

kk_uintptr clua_newfunction(struct lua_State *L);
