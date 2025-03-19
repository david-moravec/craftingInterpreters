#include "chunk.h"
#include "debug.h"
#include "value.h"
#include "vm.h"
#include <stdio.h>

void resetStack(VM* vm) { vm->stackTop = vm->stack; }

void push(Value** stackTop, Value value) {
  **stackTop = value;
  (*stackTop)++;
}

Value pop(Value** stackTop) {
  (*stackTop)--;
  return **stackTop;
}

void initVM(VM* vm) {
  resetStack(vm);
  vm->chunk = NULL;
  vm->ip = NULL;
}

void freeVM(VM* vm) {
  freeChunk(vm->chunk);
  initVM(vm);
}

static InterpretResult run(VM* vm) {
#define READ_BYTE() (*vm->ip++)
#define READ_CONSTANT() (vm->chunk->constants.values[READ_BYTE()])
  for (;;) {
#ifdef DEBUG_TRACE_EXECUTION
    for (Value* slot = vm->stack; slot < vm->stackTop; slot++) {
      printf("[");
      printValue(*slot);
      printf("]\n");
    }
    disassembleInstruction(vm->chunk, (int)(vm->ip - vm->chunk->code));
#endif
    byte instruction = READ_BYTE();
    switch (instruction) {
    case OP_CONSTANT: {
      Value value = READ_CONSTANT();
      push(&vm->stackTop, value);
      break;
    }
    case OP_RETURN: {
      printValue(pop(&vm->stackTop));
      printf("\n");
      return INTERPRET_OK;
    }
    }
  }
#undef READ_CONSTANT
#undef READ_BYTE
}

InterpretResult interpret(Chunk* chunk, VM* vm) {
  vm->chunk = chunk;
  vm->ip = vm->chunk->code;
  return run(vm);
}
