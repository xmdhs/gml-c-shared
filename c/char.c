#include "char.h"

char **newchar(long long i)
{
    char **c = malloc(i * sizeof(char *));
    return c;
}

void setchar(char **c, long long index, char *s)
{
    c[index] = s;
}

char *getchar(char **c, long long index)
{
    char *s = c[index];
    return s;
}

void freechar(char **c, long long len)
{
    for (long long i = 0; i < len; i++)
    {
        free(c[i]);
    }
    free(c);
}