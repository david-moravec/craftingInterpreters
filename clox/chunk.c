#include "chunk.h"
#include "memory.h"
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

void initChunk(Chunk* chunk) {
  chunk->count = 0;
  chunk->capacity = 0;
  chunk->code = NULL;
}

void writeChunk(Chunk* chunk, byte b) {
  if (chunk->capacity < chunk->count + 1) {
    int old_capacity = chunk->capacity;
    chunk->capacity = GROW_CAPACITY(old_capacity);
    chunk->code = GROW_ARRAY(byte, chunk->code, old_capacity, chunk->capacity);
  }
  chunk->code[chunk->count] = b;
  chunk->count++;
}

void freeChunk(Chunk* chunk) {
  FREE_ARRAY(byte, chunk->code, chunk->capacity);
  initChunk(chunk);
}
