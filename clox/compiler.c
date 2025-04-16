#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#include "chunk.h"
#include "common.h"
#include "compiler.h"
#include "object.h"
#include "scanner.h"
#include "value.h"
#include "vm.h"

#ifdef DEBUG_PRINT_CODE
#include "debug.h"
#endif

typedef struct {
  Token current;
  Token previous;
  bool hadError;
  bool panicMode;
} Parser;

typedef enum {
  PREC_NONE,
  PREC_ASSIGNMENT, // =
  PREC_OR,         // or
  PREC_AND,        // and
  PREC_EQUALITY,   // == !=
  PREC_COMPARISON, // < > <= >=
  PREC_TERM,       // + -
  PREC_FACTOR,     // * /
  PREC_UNARY,      // ! -
  PREC_CALL,       // . ()
  PREC_PRIMARY
} Precedence;

typedef void (*ParseFn)(VM*, Parser*, Scanner*);

typedef struct {
  ParseFn prefix;
  ParseFn infix;
  Precedence precedence;
} ParseRule;

static void expression(VM* vm, Parser* parser, Scanner* scanner);
static void statement(VM* vm, Parser* parser, Scanner* scanner);
static void declaration(VM* vm, Parser* parser, Scanner* scanner);
static ParseRule* getRule(TokenType type);
static void parsePrecedence(VM* vm, Parser* parser, Scanner* scanner,
                            Precedence precedence);

Chunk* compilingChunk;

static Chunk* currentChunk() { return compilingChunk; }

static void errorAt(Parser* parser, const Token* token, const char* message) {
  if (parser->panicMode)
    return;
  parser->panicMode = true;
  fprintf(stderr, "[line %d] Error", token->line);

  if (token->type == TOKEN_EOF) {
    fprintf(stderr, " at end");
  } else if (token->type == TOKEN_ERROR) {
    // Nothing.
  } else {
    fprintf(stderr, " at '%.*s'", token->length, token->start);
  }

  fprintf(stderr, ": %s\n", message);
  parser->hadError = true;
}

static void error(Parser* parser, const char* message) {
  errorAt(parser, &parser->previous, message);
}

static void errorAtCurrent(Parser* parser, const char* message) {
  errorAt(parser, &parser->current, message);
}

static void advance(Parser* parser, Scanner* scanner) {
  parser->previous = parser->current;
  for (;;) {
    parser->current = scanToken(scanner);
    if (parser->current.type != TOKEN_ERROR)
      break;

    errorAtCurrent(parser, parser->current.start);
  }
}

static void consume(Parser* parser, Scanner* scanner, TokenType expected,
                    const char* message) {
  if (parser->current.type == expected) {
    advance(parser, scanner);
    return;
  }
  errorAtCurrent(parser, message);
}

static bool check(Parser* parser, TokenType type) {
  return parser->current.type == type;
}

static bool match(Parser* parser, Scanner* scanner, TokenType type) {
  if (!check(parser, type))
    return false;
  advance(parser, scanner);
  return true;
}

static void emitByte(const Parser* parser, const byte b) {
  writeChunk(currentChunk(), b, parser->previous.line);
}

static void emitBytes(const Parser* parser, const byte b1, const byte b2) {
  emitByte(parser, b1);
  emitByte(parser, b2);
}

static void emitReturn(Parser* parser) { emitByte(parser, OP_RETURN); }

static void endCompiler(Parser* parser) {
  emitReturn(parser);
#ifdef DEBUG_PRINT_CODE
  if (parser->hadError) {
    disassembleChunk(currentChunk(), "code");
  }
#endif
}

static void expression(VM* vm, Parser* parser, Scanner* scanner) {
  parsePrecedence(vm, parser, scanner, PREC_ASSIGNMENT);
}

static void printStatement(VM* vm, Parser* parser, Scanner* scanner) {
  expression(vm, parser, scanner);
  consume(parser, scanner, TOKEN_SEMICOLON, "Expected ';' after expression.");
  emitByte(parser, OP_PRINT);
}

static void synchronize(Parser* parser, Scanner* scanner) {
  parser->panicMode = false;

  while (parser->current.type != TOKEN_EOF) {
    if (parser->previous.type == TOKEN_SEMICOLON)
      return;
    switch (parser->current.type) {
    case TOKEN_CLASS:
    case TOKEN_FUN:
    case TOKEN_PRINT:
    case TOKEN_VAR:
    case TOKEN_WHILE:
    case TOKEN_FOR:
    case TOKEN_IF:
    case TOKEN_RETURN:
      return;
    default:;
    }

    advance(parser, scanner);
  }
}

static void expressionStatement(VM* vm, Parser* parser, Scanner* scanner) {
  expression(vm, parser, scanner);
  consume(parser, scanner, TOKEN_SEMICOLON, "Expected ';' after expression.");
  emitByte(parser, OP_POP);
}

static void binary(VM* vm, Parser* parser, Scanner* scanner) {
  TokenType operatorType = parser->previous.type;
  ParseRule* rule = getRule(operatorType);
  parsePrecedence(vm, parser, scanner, (Precedence)(rule->precedence + 1));

  switch (operatorType) {
  case TOKEN_BANG_EQUAL:
    emitBytes(parser, OP_EQUAL, OP_NOT);
    break;
  case TOKEN_EQUAL_EQUAL:
    emitByte(parser, OP_EQUAL);
    break;
  case TOKEN_GREATER:
    emitByte(parser, OP_GREATER);
    break;
  case TOKEN_GREATER_EQUAL:
    emitBytes(parser, OP_LESS, OP_NOT);
    break;
  case TOKEN_LESS:
    emitByte(parser, OP_LESS);
    break;
  case TOKEN_LESS_EQUAL:
    emitBytes(parser, OP_GREATER, OP_NOT);
    break;
  case TOKEN_PLUS:
    emitByte(parser, OP_ADD);
    break;
  case TOKEN_MINUS:
    emitByte(parser, OP_SUBTRACT);
    break;
  case TOKEN_STAR:
    emitByte(parser, OP_MULTIPLY);
    break;
  case TOKEN_SLASH:
    emitByte(parser, OP_DIVIDE);
    break;
  default:
    return;
  }
}

static void grouping(VM* vm, Parser* parser, Scanner* scanner) {
  expression(vm, parser, scanner);
  consume(parser, scanner, TOKEN_RIGHT_PAREN, "Expect ')' after expression.");
}

static void unary(VM* vm, Parser* parser, Scanner* scanner) {
  TokenType operatorType = parser->previous.type;

  parsePrecedence(vm, parser, scanner, PREC_UNARY);

  switch (operatorType) {
  case TOKEN_BANG:
    emitByte(parser, OP_NOT);
    break;
  case TOKEN_MINUS:
    emitByte(parser, OP_NEGATE);
    break;
  default:
    return;
  }
}

static byte makeConstant(Parser* parser, Value value) {
  int constant = addConstant(currentChunk(), value);
  if (constant > UINT8_MAX) {
    error(parser, "Too many constants in one chunk.");
    return 0;
  }

  return (byte)constant;
}

static void emitConstant(Parser* parser, Value value) {
  emitBytes(parser, OP_CONSTANT, makeConstant(parser, value));
}

static void number(VM* vm, Parser* parser, Scanner* scanner) {
  double value = strtod(parser->previous.start, NULL);
  emitConstant(parser, NUMBER_VAL(value));
}

static void literal(VM* vm, Parser* parser, Scanner* scanner) {
  switch (parser->previous.type) {
  case TOKEN_FALSE:
    emitByte(parser, OP_FALSE);
    break;
  case TOKEN_TRUE:
    emitByte(parser, OP_TRUE);
    break;
  case TOKEN_NIL:
    emitByte(parser, OP_NIL);
    break;
  default:
    return; // unreachable
  }
}

static void string(VM* vm, Parser* parser, Scanner* scanner) {
  emitConstant(parser, OBJ_VAL(copyString(vm, parser->previous.start + 1,
                                          parser->previous.length - 2)));
}

ParseRule rules[] = {
    [TOKEN_LEFT_PAREN] = {grouping, NULL, PREC_NONE},
    [TOKEN_RIGHT_PAREN] = {NULL, NULL, PREC_NONE},
    [TOKEN_LEFT_BRACE] = {NULL, NULL, PREC_NONE},
    [TOKEN_RIGHT_BRACE] = {NULL, NULL, PREC_NONE},
    [TOKEN_COMMA] = {NULL, NULL, PREC_NONE},
    [TOKEN_DOT] = {NULL, NULL, PREC_NONE},
    [TOKEN_MINUS] = {unary, binary, PREC_TERM},
    [TOKEN_PLUS] = {NULL, binary, PREC_TERM},
    [TOKEN_SEMICOLON] = {NULL, NULL, PREC_NONE},
    [TOKEN_SLASH] = {NULL, binary, PREC_FACTOR},
    [TOKEN_STAR] = {NULL, binary, PREC_FACTOR},
    [TOKEN_BANG] = {unary, NULL, PREC_NONE},
    [TOKEN_BANG_EQUAL] = {NULL, binary, PREC_EQUALITY},
    [TOKEN_EQUAL] = {NULL, NULL, PREC_NONE},
    [TOKEN_EQUAL_EQUAL] = {NULL, binary, PREC_EQUALITY},
    [TOKEN_GREATER] = {NULL, binary, PREC_EQUALITY},
    [TOKEN_GREATER_EQUAL] = {NULL, binary, PREC_EQUALITY},
    [TOKEN_LESS] = {NULL, binary, PREC_EQUALITY},
    [TOKEN_LESS_EQUAL] = {NULL, binary, PREC_EQUALITY},
    [TOKEN_IDENTIFIER] = {NULL, NULL, PREC_NONE},
    [TOKEN_STRING] = {string, NULL, PREC_NONE},
    [TOKEN_NUMBER] = {number, NULL, PREC_NONE},
    [TOKEN_AND] = {NULL, NULL, PREC_NONE},
    [TOKEN_CLASS] = {NULL, NULL, PREC_NONE},
    [TOKEN_ELSE] = {NULL, NULL, PREC_NONE},
    [TOKEN_FALSE] = {literal, NULL, PREC_NONE},
    [TOKEN_FOR] = {NULL, NULL, PREC_NONE},
    [TOKEN_FUN] = {NULL, NULL, PREC_NONE},
    [TOKEN_IF] = {NULL, NULL, PREC_NONE},
    [TOKEN_NIL] = {literal, NULL, PREC_NONE},
    [TOKEN_OR] = {NULL, NULL, PREC_NONE},
    [TOKEN_PRINT] = {NULL, NULL, PREC_NONE},
    [TOKEN_RETURN] = {NULL, NULL, PREC_NONE},
    [TOKEN_SUPER] = {NULL, NULL, PREC_NONE},
    [TOKEN_THIS] = {NULL, NULL, PREC_NONE},
    [TOKEN_TRUE] = {literal, NULL, PREC_NONE},
    [TOKEN_VAR] = {NULL, NULL, PREC_NONE},
    [TOKEN_WHILE] = {NULL, NULL, PREC_NONE},
    [TOKEN_ERROR] = {NULL, NULL, PREC_NONE},
    [TOKEN_EOF] = {NULL, NULL, PREC_NONE},
};

static ParseRule* getRule(TokenType type) { return &rules[type]; }

static void parsePrecedence(VM* vm, Parser* parser, Scanner* scanner,
                            Precedence precedence) {
  advance(parser, scanner);
  ParseFn prefix = getRule(parser->previous.type)->prefix;
  if (prefix == NULL) {
    error(parser, "Expected expression.");
    return;
  }

  prefix(vm, parser, scanner);

  while (precedence <= getRule(parser->current.type)->precedence) {
    advance(parser, scanner);
    ParseFn infix = getRule(parser->previous.type)->infix;
    infix(vm, parser, scanner);
  }
}

static void declaration(VM* vm, Parser* parser, Scanner* scanner) {
  statement(vm, parser, scanner);

  if (parser->panicMode)
    synchronize(parser, scanner);
}

static void statement(VM* vm, Parser* parser, Scanner* scanner) {
  if (match(parser, scanner, TOKEN_PRINT)) {
    printStatement(vm, parser, scanner);
  } else {
    expressionStatement(vm, parser, scanner);
  }
}

bool compile(VM* vm, const char* source, Chunk* chunk) {
  Scanner scanner;
  Parser parser;
  initScanner(source, &scanner);
  compilingChunk = chunk;

  parser.hadError = false;
  parser.panicMode = false;

  advance(&parser, &scanner);

  while (!match(&parser, &scanner, TOKEN_EOF)) {
    declaration(vm, &parser, &scanner);
  }

  endCompiler(&parser);

  return !parser.hadError;
}
