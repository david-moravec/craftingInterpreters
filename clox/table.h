#ifndef clox_table_h
#define clox_table_h

#include <stdbool.h>
#include <stdlib.h>

#include "common.h"
#include "object.h"
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
Entry* findEntry(Entry* entries, int capacity, ObjString* key);
bool tableGet(Table* table, ObjString* key, Value* value);
bool tableSet(Table* table, ObjString* key, Value value);
bool tableDelete(Table* table, ObjString* key);
void tableAddAll(Table* from, Table* to);
ObjString* tableFindString(Table* table, char* chars, int length, uint32_t hash);

#endif
