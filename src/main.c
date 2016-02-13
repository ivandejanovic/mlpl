/*
 ============================================================================
 Name        : MLPL
 Author      : Ivan Dejanovic
 Version     : 0.1
 Copyright   : MIT License
 Description : MLPL Interpreter
 ============================================================================
 */

#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include "vector.h"
#include "types.h"
#include "config.h"
#include "parse.h"
#include "lex.h"

#define FILENAMELENGTH 120

int main(int argc, char* argv[])
{
    /* check arguments number*/
    if (argc < 2 || argc > 3)
    {
        fprintf (stderr, "Usage: <codefilename> [configurationfilename].\n");
        return EXIT_FAILURE;
    }
    char configFilename[FILENAMELENGTH];
    char codeFilename[FILENAMELENGTH];
    ReservedWord* reserved = NULL;

    /* create reserved words*/
    if(argc == 3)
    {
        strncpy(configFilename, argv[2], FILENAMELENGTH);
        FILE* configFile = fopen(configFilename, "r");
        if (configFile == NULL)
        {
            fprintf(stderr, "Error opening configuration file.\n");
            return EXIT_FAILURE;
        }
        reserved = getConfigReserved(configFile);
        fclose(configFile);
    }
    else
    {
        reserved = getDefaultReserved();
    }

    strncpy(codeFilename, argv[1], FILENAMELENGTH);
    FILE* codeFile = fopen(codeFilename, "r");
    if(codeFile == NULL)
    {
        fprintf(stderr, "Error opening code file.\n");
        return EXIT_FAILURE;
    }

    vector* tokens = parse(codeFile, reserved);

    /* close source file */
    fclose(codeFile);
    /* release reserved words from memory */
    free(reserved);

    lex(tokens);

    return EXIT_SUCCESS;
}
