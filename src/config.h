/*
 ============================================================================
 Name        : config
 Author      : Ivan Dejanovic
 Copyright   : MIT License
 Description : configuration files handler header file
 ============================================================================
 */

#ifndef _CONFIG_H_
#define _CONFIG_H_

ReservedWord* getDefaultReserved(void);
ReservedWord* getConfigReserved(FILE* config);

#endif /* _CONFIG_H_ */
