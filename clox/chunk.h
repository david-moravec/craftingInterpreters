#ifndef clox_chunk_h
#define clox_chunk_h

#include "common.h"
#include <stdint.h>

typedef enum { OP_RETURN } OpCode;
typedef uint8_t byte;

typedef struct {
  int count;
  int capacity;
  byte* code;
} Chunk;

void initChunk(Chunk* chunk);
void writeChunk(Chunk* chunk, byte byte);
void freeChunk(Chunk* chunk);

#endif
