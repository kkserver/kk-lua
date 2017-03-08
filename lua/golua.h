#include <stdint.h>

typedef void * uintptr;

uintptr clua_newobject(struct lua_State * L);

uintptr clua_newfunction(struct lua_State *L);
