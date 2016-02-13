/*
 ============================================================================
 Name        : vector
 Author      : Ivan Dejanovic
 Copyright   : MIT License
 Description : vector header file
 Disclaimer  : Initial implementation was taken from http://eddmann.com/posts/implementing-a-dynamic-vector-array-in-c/
 ============================================================================
 */

#ifndef _VECTOR_H_
#define _VECTOR_H_

#define VECTOR_INIT_CAPACITY 8

#define VECTOR_INIT(vec, clear) vector_init(&vec, clear)
#define VECTOR_ADD(vec, item) vector_add(&vec, (void*) item)
#define VECTOR_SET(vec, id, item) vector_set(&vec, id, (void*) item)
#define VECTOR_GET(vec, type, id) (type) vector_get(&vec, id)
#define VECTOR_DELETE(vec, id) vector_delete(&vec, id)
#define VECTOR_SIZE(vec) vector_size(&vec)
#define VECTOR_FREE(vec) vector_free(&vec)

typedef void(*clear_ptr)(void*);

typedef struct vector {
    void** items;
    long capacity;
    long size;
    clear_ptr clear;
} vector;

void vector_init(vector*, clear_ptr);
long vector_size(vector*);
void vector_add(vector*, void*);
void vector_set(vector*, long, void*);
void* vector_get(vector*, long);
void vector_delete(vector*, long);
void vector_free(vector*);

#endif /* _VECTOR_H_ */
