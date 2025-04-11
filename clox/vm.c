#include <stdarg.h>
#include <stdio.h>
#include <string.h>
#include <vadefs.h>

#include "chunk.h"
#include "compiler.h"
#include "debug.h"
#include "memory.h"
#include "object.h"
#include "table.h"
#include "value.h"
#include "vm.h"

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
  vm->objects = NULL;
  initTable(&vm->strings);
}

void freeVM(VM* vm) {
  freeChunk(vm->chunk);
  freeTable(&vm->strings);
  freeObjects(vm->objects);
  initVM(vm);
}

static Value peek(VM* vm, int distance) { return vm->stackTop[-1 - distance]; }

static void runtimeError(VM* vm, const char* format, ...) {
  va_list args;
  va_start(args, format);
  vfprintf(stderr, format, args);
  va_end(args);
  fputs("\n", stderr);

  size_t instuction = vm->ip - vm->chunk->code - 1;
  int line = vm->chunk->lines[instuction];
  fprintf(stderr, "[line %d] in script\n", line);
  resetStack(vm);
}

static bool isFalsey(Value value) {
  return IS_NIL(value) || (IS_BOOL(value) && !AS_BOOL(value));
}

static void concatenate(VM* vm) {
  ObjString* a = AS_STRING(pop(&vm->stackTop));
  ObjString* b = AS_STRING(pop(&vm->stackTop));

  int length = a->length + b->length;
  char* chars = ALLOCATE(char, length + 1);
  memcpy(chars, a, a->length);
  memcpy(chars + a->length, b, b->length);
  chars[length] = '\0';

  ObjString* result = takeString(vm, chars, length);
  push(&vm->stackTop, OBJ_VAL(result));
}

static InterpretResult run(VM* vm) {
#define READ_BYTE() (*vm->ip++)
#define READ_CONSTANT() (vm->chunk->constants.values[READ_BYTE()])
#define POP() pop(&vm->stackTop)
#define PUSH(value) push(&vm->stackTop, value)
#define BINARY_OP(valueType, op)                                               \
  do {                                                                         \
    if (!IS_NUMBER(peek(vm, 0)) || !IS_NUMBER(peek(vm, 1))) {                  \
      runtimeError(vm, "Operands must be numbers.");                           \
      return INTERPRET_RUNTIME_ERROR;                                          \
    }                                                                          \
    double b = AS_NUMBER(POP());                                               \
    double a = AS_NUMBER(POP());                                               \
    PUSH(valueType(a op b));                                                   \
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
    case OP_NIL: {
      PUSH(NIL_VAL);
      break;
    }
    case OP_TRUE: {
      PUSH(BOOL_VAL(true));
      break;
    }
    case OP_FALSE: {
      PUSH(BOOL_VAL(false));
      break;
    }
    case OP_EQUAL: {
      Value a = POP();
      Value b = POP();
      PUSH(BOOL_VAL(valuesEqual(a, b)));
      break;
    }
    case OP_GREATER: {
      BINARY_OP(BOOL_VAL, >);
      break;
    }
    case OP_LESS: {
      BINARY_OP(BOOL_VAL, <);
      break;
    }
    case OP_ADD: {
      if (IS_STRING(peek(vm, 0)) && IS_STRING(peek(vm, 1))) {
        concatenate(vm);
      } else if (IS_NUMBER(peek(vm, 0)) && IS_NUMBER(peek(vm, 1))) {
        double b = AS_NUMBER(POP());
        double a = AS_NUMBER(POP());
        PUSH(NUMBER_VAL(a + b));
      } else {
        runtimeError(vm, "Operands must be two numbers or two strings.");
        return INTERPRET_RUNTIME_ERROR;
      }
      break;
    }
    case OP_SUBTRACT: {
      BINARY_OP(NUMBER_VAL, -);
      break;
    }
    case OP_MULTIPLY: {
      BINARY_OP(NUMBER_VAL, *);
      break;
    }
    case OP_DIVIDE: {
      BINARY_OP(NUMBER_VAL, /);
      break;
    }
    case OP_NOT: {
      PUSH(BOOL_VAL(isFalsey(POP())));
      break;
    }
    case OP_NEGATE: {
      if (!IS_NUMBER(peek(vm, 0))) {
        runtimeError(vm, "Operand must be number.");
        return INTERPRET_RUNTIME_ERROR;
      }
      PUSH(NUMBER_VAL(-AS_NUMBER(POP())));
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
  if (!compile(vm, source, &chunk)) {
    freeChunk(&chunk);
    return INTERPRET_COMPILE_ERROR;
  };

  vm->chunk = &chunk;
  vm->ip = vm->chunk->code;

  InterpretResult result = run(vm);

  freeChunk(&chunk);
  return result;
}
