#ifndef clox_table_h
#define clox_table_h

#include <stdbool.h>
#include <stdlib.h>

#include "common.h"
#include "value.h"

typedef struct {
  ObjString* key;
  Value value;
} Entry;

typedef struct {
  int count;
  int capacity;
  Entry* entries;
} Table;

void initTable(Table* table);
void freeTable(Table* table);
Entry* findEntry(Table* table, ObjString* key);
bool tableSet(Table* table, ObjString* key, Value value);

#endif
