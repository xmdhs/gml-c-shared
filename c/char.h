#ifndef _GO_CHAR_
#define _GO_CHAR_

#include <stdlib.h>

typedef struct
{
    int code;
    char *msg;
} err;

char **newchar(long long i);

void setchar(char **c, long long index, char *s);

char *getchar(char **c, long long index);

void freechar(char **c, long long len);

void do_Fail(void *f, char *c);

void do_Ok(void *f, int i1, int i2);

void do_finish(void *f,err e);

void* GoMalloc(int i);

#endif
