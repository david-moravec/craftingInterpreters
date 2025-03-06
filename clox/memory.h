#ifndef clox_memory_h
#define clox_memory_h

#include "common.h"

#define GROW_CAPACITY(capacity) ((capacity) < 8 ? 8 : (capacity) * 2)

#define GROW_ARRAY(type, pointer, old_count, new_count)                        \
  (type*)reallocate(pointer, (old_count) * sizeof(type),                       \
                    (new_count) * sizeof(type))

#define FREE_ARRAY(type, pointer, old_count)                                   \
  (type*)reallocate(pointer, (old_count) * sizeof(type), 0)

void* reallocate(void* pointer, size_t old_size, size_t new_size);

#endif
