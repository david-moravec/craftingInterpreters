#include <stdio.h>

#include "common.h"
#include "compiler.h"
#include "scanner.h"

typedef struct {
  Token current;
  Token previous;
} Parser;

void advance(Parser* parser, Scanner* scanner) {
  parser->previous = parser->current;
  for (;;) {
    parser->current = scanToken(scanner);
    if (parser->current != TOKEN_ERROR)
      break;

    errorAtCurrent(parser->current.start);
  }
}

void compile(const char* source) {
  Scanner scanner;
  initScanner(source, &scanner);
  int line = -1;
  for (;;) {
    Token token = scanToken(&scanner);
    if (token.line != line) {
      printf("%4d ", token.line);
      line = token.line;
    } else {
      printf("    | ");
    }
    printf("%2d '%.*s'\n", token.type, token.length, token.start);

    if (token.type == TOKEN_EOF)
      break;
  }
}
