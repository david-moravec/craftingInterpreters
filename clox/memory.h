#ifndef clox_memory_h
#define clox_memory_h

#include "common.h"
#include "object.h"

#define ALLOCATE(type, count) \
  (type*)reallocate(NULL, 0, sizeof(type) * count)

#define FREE(type, pointer) reallocate(pointer, sizeof(type), 0)

#define GROW_CAPACITY(capacity) ((capacity) < 8 ? 8 : (capacity) * 2)

#define GROW_ARRAY(type, pointer, old_count, new_count)                        \
  (type*)reallocate(pointer, (old_count) * sizeof(type),                       \
                    (new_count) * sizeof(type))

#define FREE_ARRAY(type, pointer, old_count)                                   \
  (type*)reallocate(pointer, (old_count) * sizeof(type), 0)

void* reallocate(void* pointer, size_t old_size, size_t new_size);
void freeObjects(Obj* objects);

#endif
