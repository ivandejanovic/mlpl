/*
 ============================================================================
 Name        : config
 Author      : Ivan Dejanovic
 Copyright   : MIT License
 Description : configuration files handler source file
 ============================================================================
 */

#include <stdlib.h>
#include <stdio.h>

#include "types.h"
#include "config.h"

ReservedWord* fillDefault(ReservedWord* reserved)
{
    /* Add default reserved words */
    reserved[0].str = "if";
    reserved[0].tok = IF;

    reserved[1].str = "then";
    reserved[1].tok = THEN;

    reserved[2].str = "else";
    reserved[2].tok = ELSE;

    reserved[3].str = "end";
    reserved[3].tok = END;

    reserved[4].str = "repeat";
    reserved[4].tok = REPEAT;

    reserved[5].str = "until";
    reserved[5].tok = UNTIL;

    reserved[6].str = "read";
    reserved[6].tok = READ;

    reserved[7].str = "write";
    reserved[7].tok = WRITE;

    return reserved;
}

ReservedWord* getDefaultReserved(void)
{
    ReservedWord* reserved = calloc(MAXRESERVED, sizeof(ReservedWord));

    if(reserved == NULL)
    {
        fprintf(stderr, "Failure to allocate memory for reserved words.\n");
        exit(EXIT_FAILURE);
    }

    return fillDefault(reserved);
}

ReservedWord*
getConfigReserved(FILE* config)
{
    return getDefaultReserved();
}
