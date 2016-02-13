/*
 ============================================================================
 Name        : parse
 Author      : Ivan Dejanovic
 Copyright   : MIT License
 Description : parser source file
 ============================================================================
 */

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <ctype.h>

#include "vector.h"
#include "types.h"
#include "parse.h"

/* states in parser */
typedef enum stateType
{
    START, INASSIGN, INCOMMENT, INNUM, INID, DONE
} StateType;

static TokenType reservedLookup(char* s, ReservedWord* reserved)
{
    int i;
    for(i = 0; i < MAXRESERVED; i++)
    {
        if(!strcmp(s, reserved[i].str))
        {
            return reserved[i].tok;
        }
    }
    return ID;
}

static Token* getToken(FILE* source, ReservedWord* reserved)
{
    /* lexeme of identifier or reserved word */
    char tokenString[MAXTOKENLEN + 1];
    /* index for storing into tokenString */
    int tokenStringIndex = 0;
    /* holds current token to be returned */
    TokenType currentToken;
    /* current state - always begins at START */
    StateType state = START;
    /* flag to indicate save to tokenString */
    int save;
    /* ungetchar return code */
    int ret_code;

    while(state != DONE)
    {
        int c = fgetc(source);
        save = TRUE;
        switch(state)
        {
            case START:
                if(isdigit(c))
                {
                    state = INNUM;
                }
                else if(isalpha(c))
                {
                    state = INID;
                }
                else if(c == ':')
                {
                    state = INASSIGN;
                }
                else if((c == ' ') || (c == '\t') || (c == '\n'))
                {
                    save = FALSE;
                }
                else if(c == '#')
                {
                    save = FALSE;
                    state = INCOMMENT;
                }
                else
                {
                    state = DONE;
                    switch(c)
                    {
                        case EOF:
                            save = FALSE;
                            currentToken = ENDFILE;
                            break;
                        case '=':
                            currentToken = EQ;
                            break;
                        case '<':
                            currentToken = LT;
                            break;
                        case '+':
                            currentToken = PLUS;
                            break;
                        case '-':
                            currentToken = MINUS;
                            break;
                        case '*':
                            currentToken = TIMES;
                            break;
                        case '/':
                            currentToken = OVER;
                            break;
                        case '(':
                            currentToken = LPAREN;
                            break;
                        case ')':
                            currentToken = RPAREN;
                            break;
                        case ';':
                            currentToken = SEMI;
                            break;
                        default:
                            currentToken = ERROR;
                            break;
                    }
                }
                break;
            case INCOMMENT:
                save = FALSE;
                if(c == EOF)
                {
                    state = DONE;
                    currentToken = ENDFILE;
                }
                else if(c == '#')
                {
                    state = START;
                }
                break;
            case INASSIGN:
                state = DONE;
                if(c == '=')
                {
                    currentToken = ASSIGN;
                }
                else
                {
                    /* backup in the input */
                    ret_code = ungetc(c, source);
                    if(ret_code == EOF)
                    {
                        if(ferror(source))
                        {
                            fprintf(stderr, "ungetc() failed in file %s at line # %d\n", __FILE__, __LINE__ - 4);
                            exit(EXIT_FAILURE);
                        }
                    }
                    save = FALSE;
                    currentToken = ERROR;
                }
                break;
            case INNUM:
                if(!isdigit(c))
                {
                    /* backup in the input */
                    ungetc(c, source);
                    if(ret_code == EOF)
                    {
                        if(ferror(source))
                        {
                            fprintf(stderr, "ungetc() failed in file %s at line # %d\n", __FILE__, __LINE__ - 4);
                            exit(EXIT_FAILURE);
                        }
                    }
                    save = FALSE;
                    state = DONE;
                    currentToken = NUM;
                }
                break;
            case INID:
                if(!isalpha(c))
                {
                    /* backup in the input */
                    ungetc(c, source);
                    if(ret_code == EOF)
                    {
                        if(ferror(source))
                        {
                            fprintf(stderr, "ungetc() failed in file %s at line # %d\n", __FILE__, __LINE__ - 4);
                            exit(EXIT_FAILURE);
                        }
                    }
                    save = FALSE;
                    state = DONE;
                    currentToken = ID;
                }
                break;
            case DONE:
            default:
                /* should never happen */
                fprintf(stderr, "Scanner Bug: state= %d\n", state);
                state = DONE;
                currentToken = ERROR;
                break;
        }
        if((save) &&(tokenStringIndex <= MAXTOKENLEN))
        {
            tokenString[tokenStringIndex++] = (char) c;
        }
        if(state == DONE)
        {
            tokenString[tokenStringIndex] = '\0';
            if(currentToken == ID)
            {
                currentToken = reservedLookup(tokenString, reserved);
            }
        }
    }
    Token* token = malloc(sizeof(Token));
    token->type = currentToken;
    strncpy(token->tokenString, tokenString, MAXTOKENLEN + 1);
    return token;
}

vector* parse(FILE* source, ReservedWord* reserved)
{
    vector* tokens = malloc(sizeof(vector));
    vector_init(tokens, NULL);

    Token* token = NULL;

    do
    {
        token = getToken(source, reserved);
        vector_add(tokens, (void*)token);
    }
    while(token->type != ENDFILE);

    return tokens;
}
