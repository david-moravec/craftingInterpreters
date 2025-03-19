#include "chunk.h"
#include "common.h"
#include "debug.h"
#include "vm.h"

int main(int argc, const char* argv[]) {
  VM vm;
  initVM(&vm);
  Chunk chunk;
  initChunk(&chunk);
  int constant = addConstant(&chunk, 2.5);
  writeChunk(&chunk, OP_CONSTANT, 1);
  writeChunk(&chunk, constant, 1);
  writeChunk(&chunk, OP_RETURN, 1);

  interpret(&chunk, &vm);

  freeVM(&vm);
  freeChunk(&chunk);
  return 0;
}
