#ifndef clox_compiler_h
#define clox_compiler_h

#include "chunk.h"
#include "vm.h"

bool compile(VM* vm, const char* source, Chunk* chunk);

#endif
