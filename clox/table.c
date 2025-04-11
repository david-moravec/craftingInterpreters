#include <stdbool.h>
#include <stdlib.h>

#include "memory.h"
#include "table.h"

void initTable(Table* table) {
  table->count = 0;
  table->capacity = 0;
  table->entries = NULL;
}

void freeTable(Table* table) {
  FREE_ARRAY(Entry, table->entries, table->capacity);
  initTable(table);
}

bool tableSet(Table* table, ObjString* key, Value value) {
  Entry* entry = findEntry(table, key);
  bool isNewKey = entry == NULL;

  if (isNewKey) {
    table->count++;
  }
  entry->key = key;
  entry->value = value;
  return isNewKey;
}
