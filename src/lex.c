/*
 ============================================================================
 Name        : lex
 Author      : Ivan Dejanovic
 Copyright   : MIT License
 Description : lexer source file
 ============================================================================
 */

#include <stdlib.h>
#include <stdio.h>

#include "vector.h"
#include "types.h"
#include "lex.h"

/* Procedure printToken prints a token and its lexeme to the listing file */
void printToken(TokenType token, const char* tokenString)
{
    switch(token)
    {
        case IF:
        case THEN:
        case ELSE:
        case END:
        case REPEAT:
        case UNTIL:
        case READ:
        case WRITE:
            fprintf(stdout, "reserved word: %s\n", tokenString);
            break;
        case ASSIGN:
            fprintf(stdout, ":=\n");
            break;
        case LT:
            fprintf(stdout, "<\n");
            break;
        case EQ:
            fprintf(stdout, "=\n");
            break;
        case LPAREN:
            fprintf(stdout, "(\n");
            break;
        case RPAREN:
            fprintf(stdout, ")\n");
            break;
        case SEMI:
            fprintf(stdout, ";\n");
            break;
        case PLUS:
            fprintf(stdout, "+\n");
            break;
        case MINUS:
            fprintf(stdout, "-\n");
            break;
        case TIMES:
            fprintf(stdout, "*\n");
            break;
        case OVER:
            fprintf(stdout, "/\n");
            break;
        case ENDFILE:
            fprintf(stdout, "EOF\n");
            break;
        case NUM:
            fprintf(stdout, "NUM, val= %s\n", tokenString);
            break;
        case ID:
            fprintf(stdout, "ID, name= %s\n", tokenString);
            break;
        case ERROR:
            fprintf(stdout, "ERROR: %s\n", tokenString);
            break;
        default: /* should never happen */
            fprintf(stdout, "Unknown token: %d\n", token);
    }
}

TreeNode* lex(vector* tokens)
{
    long index, size;
    size = vector_size(tokens);
    Token* token;

    for (index = 0; index < size; ++index) {
        token = (Token*)vector_get(tokens, index);
        printToken(token->type, token->tokenString);
    }

    return NULL;
}
