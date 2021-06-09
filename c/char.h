#ifndef _GO_CHAR_
#define _GO_CHAR_

#include <stdlib.h>

char **newchar(long long i);

void setchar(char **c, long long index, char *s);

char *getchar(char **c, long long index);

void freechar(char **c,long long len);

#endif
