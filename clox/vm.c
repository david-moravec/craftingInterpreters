#include "chunk.h"
#include "debug.h"
#include "value.h"
#include "compiler.h"
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
#define POP() pop(&vm->stackTop)
#define PUSH(value) push(&vm->stackTop, value)
#define BINARY_OP(op)                                                          \
  do {                                                                         \
    double b = POP();                                                          \
    double a = POP();                                                          \
    PUSH(a op b);                                                              \
  } while (false)

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
    case OP_ADD: {
      BINARY_OP(+);
      break;
    }
    case OP_SUBTRACT: {
      BINARY_OP(-);
      break;
    }
    case OP_MULTIPLY: {
      BINARY_OP(*);
      break;
    }
    case OP_DIVIDE: {
      BINARY_OP(/);
      break;
    }
    case OP_NEGATE: {
      PUSH(-POP());
      break;
    }
    case OP_RETURN: {
      printValue(POP());
      printf("\n");
      return INTERPRET_OK;
    }
    }
  }
#undef BINARY_OP
#undef PUSH
#undef POP
#undef READ_CONSTANT
#undef READ_BYTE
}

InterpretResult interpret(const char* source, VM* vm) {
  Chunk chunk;
  initChunk(&chunk);
  if (!compile(source, &chunk)) {
    freeChunk(&chunk);
    return INTERPRET_COMPILE_ERROR;
  };

  vm->chunk = &chunk;
  vm->ip = vm->chunk->code;

  InterpretResult result = run(vm);

  freeChunk(&chunk);
  return result;
}
