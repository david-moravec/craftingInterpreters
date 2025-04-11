#ifndef clox_vm_h
#define clox_vm_h

#include "chunk.h"
#include "value.h"
#include "table.h"

#define STACK_MAX 256

typedef struct {
  Chunk* chunk;
  uint8_t* ip;
  Value stack[STACK_MAX];
  Value* stackTop;
  Table strings;
  Obj* objects;
} VM;

typedef enum {
  INTERPRET_OK,
  INTERPRET_COMPILE_ERROR,
  INTERPRET_RUNTIME_ERROR,
} InterpretResult;

void push(Value** stackTop, Value value);
Value pop(Value** stackTop);
void initVM(VM* vm);
void freeVM(VM* vm);
InterpretResult interpret(const char* source, VM* vm);

#endif
