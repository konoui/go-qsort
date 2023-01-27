// +build ignore

#include <stdio.h>
#include <errno.h>
#include <stdlib.h>

#define SWAP(a, b, count, size, tmp) \
    {                                \
        count = size;                \
        do                           \
        {                            \
            tmp = *a;                \
            *a++ = *b;               \
            *b++ = tmp;              \
        } while (--count);           \
    }

/* Copy one block of size size to another. */
#define COPY(a, b, count, size, tmp1, tmp2) \
    {                                       \
        count = size;                       \
        tmp1 = a;                           \
        tmp2 = b;                           \
        do                                  \
        {                                   \
            *tmp1++ = *tmp2++;              \
        } while (--count);                  \
    }

int myheapsort(vbase, nmemb, size, compar) void *vbase;
size_t nmemb, size;
int (*compar)(const void *, const void *);
{
    fprintf(stderr, "myheapsort is called\n");
    printf("myheapsort is called\n");
    size_t cnt,
        i, j, l;
    char tmp, *tmp1, *tmp2;
    char *base, *k, *p, *t;

    if (nmemb <= 1)
        return (0);

    if (!size)
    {
        errno = EINVAL;
        fprintf(stderr, "zero size\n");
        return (-1);
    }

    if ((k = malloc(size)) == NULL)
    {
        fprintf(stderr, "null mallock\n");
        return (-1);
    }
    /*
     * Items are numbered from 1 to nmemb, so offset from size bytes
     * below the starting address.
     */
    base = (char *)vbase - size;

    for (l = nmemb / 2 + 1; --l;)
        for (i = l; (j = i * 2) <= nmemb; i = j)
        {
            p = base + j * size;
            if (j < nmemb && compar(p, p + size) < 0)
            {
                p += size;
                ++j;
            }
            t = base + i * size;
            if (compar(p, t) <= 0)
                break;
            SWAP(t, p, cnt, size, tmp);
        }

    while (nmemb > 1)
    {
        COPY(k, base + nmemb * size, cnt, size, tmp1, tmp2);
        COPY(base + nmemb * size, base + size, cnt, size, tmp1, tmp2);
        --nmemb;

        for (i = 1; (j = i * 2) <= nmemb; i = j)
        {
            p = base + j * size;
            if (j < nmemb && compar(p, p + size) < 0)
            {
                p += size;
                ++j;
            }
            t = base + i * size;
            COPY(t, p, cnt, size, tmp1, tmp2);
        }

        for (;;)
        {
            j = i;
            i = j / 2;
            p = base + j * size;
            t = base + i * size;
            if (j == 1 || compar(k, t) < 0)
            {
                COPY(p, k, cnt, size, tmp1, tmp2);
                break;
            }
            COPY(p, t, cnt, size, tmp1, tmp2);
        }
    }
    return 0;
}
