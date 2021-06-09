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

void do_Fail(void *f, char *c)
{
    void (*Fail)(char *);
    Fail = f;
    Fail(c);
}

void do_Ok(void *f, int i1, int i2)
{
    void (*ok)(int, int);
    ok = f;
    ok(i1, i2);
}