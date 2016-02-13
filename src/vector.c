/*
 ============================================================================
 Name        : vector
 Author      : Ivan Dejanovic
 Copyright   : MIT License
 Description : vector source file
 Disclaimer  : Initial implementation was taken from http://eddmann.com/posts/implementing-a-dynamic-vector-array-in-c/
 ============================================================================
 */

#include <stdio.h>
#include <stdlib.h>

#include "vector.h"

static void vector_resize(vector* v, long capacity)
{
    void** items = realloc(v->items, sizeof(void*) * capacity);
    if (items) {
        v->items = items;
        v->capacity = capacity;
    }
}

void vector_init(vector* v, clear_ptr clear)
{
    v->capacity = VECTOR_INIT_CAPACITY;
    v->size = 0;
    v->items = malloc(sizeof(void*) * v->capacity);
    v->clear = clear;
}

long vector_size(vector* v)
{
    return v->size;
}

void vector_add(vector* v, void* item)
{
    if (v->capacity == v->size)
    {
        vector_resize(v, v->capacity * 2);
    }
    v->items[v->size++] = item;
}

void vector_set(vector* v, long index, void* item)
{
    if (index >= 0 && index < v->size)
    {
        if (v->clear != NULL)
        {
            v->clear(v->items[index]);
        }
        v->items[index] = item;
    }
}

void* vector_get(vector* v, long index)
{
    if (index >= 0 && index < v->size)
    {
        return v->items[index];
    }
    return NULL;
}

void vector_delete(vector* v, long index)
{
    if (index < 0 || index >= v->size)
    {
        return;
    }

    if(v->clear != NULL)
    {
        v->clear(v->items[index]);
    }
    v->items[index] = NULL;

    for (long i = index; i < v->size - 1; ++i) {
        v->items[i] = v->items[i + 1];
    }

    v->items[v->size - 1] = NULL;

    v->size--;

    if (v->size > 0 && v->size == v->capacity / 4)
        vector_resize(v, v->capacity / 2);
}

void vector_free(vector* v)
{
    for (long i = 0; i < v->size; ++i)
    {
        v->clear(v->items[i]);
    }
    free(v->items);
}

