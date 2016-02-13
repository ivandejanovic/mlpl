/*
 ============================================================================
 Name        : types
 Author      : Ivan Dejanovic
 Version     : 0.1
 Copyright   : MIT License
 Description : types definition header file
 ============================================================================
 */

#ifndef _TYPES_H_
#define _TYPES_H_

#ifndef FALSE
#define FALSE 0
#endif

#ifndef TRUE
#define TRUE 1
#endif

/* MAXRESERVED = the number of reserved words */
#define MAXRESERVED 8

/* MAXTOKENLEN is the maximum size of a token */
#define MAXTOKENLEN 40

typedef enum tokenTypes
{
    ENDFILE, ERROR,
    /* reserved words */
    IF, THEN, ELSE, END, REPEAT, UNTIL, READ, WRITE,
    /* multicharacter tokens */
    ID, NUM,
    /* special symbols */
    ASSIGN, EQ, LT, PLUS, MINUS, TIMES, OVER, LPAREN, RPAREN, SEMI
} TokenType;

typedef struct reservedWord
{
    char* str;
    TokenType tok;
} ReservedWord;

typedef struct token {
    TokenType type;
    char tokenString[MAXTOKENLEN + 1];
} Token;

typedef enum
{
    StmtK, ExpK
} NodeKind;

typedef enum
{
    IfK, RepeatK, AssignK, ReadK, WriteK
} StmtKind;

typedef enum
{
    OpK, ConstK, IdK
} ExpKind;

/* ExpType is used for type checking */
typedef enum
{
    Void, Integer, Boolean
} ExpType;

#define MAXCHILDREN 3

typedef struct treeNode
{
    struct treeNode * child[MAXCHILDREN];
    struct treeNode * sibling;
    int lineno;
    NodeKind nodekind;
    union
    {
        StmtKind stmt;
        ExpKind exp;
    } kind;
    union
    {
        TokenType op;
        int val;
        char * name;
    } attr;
    ExpType type; /* for type checking of exps */
} TreeNode;

#endif /* _TYPES_H_ */
