// +build ignore

#include <stdio.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>
#include <mach/mach.h>

int myheapsort(void *__base, size_t __nel, size_t __width,
               int (*_Nonnull __compar)(const void *, const void *));

typedef int cmp_t(const void *, const void *);
static inline char *med3(char *, char *, char *, cmp_t *, void *);
static inline void swapfunc(char *, char *, size_t, int, int);

#define MIN(a, b) ((a) < (b) ? a : b)

/*
 * Qsort routine from Bentley & McIlroy's "Engineering a Sort Function".
 */
#define swapcode(TYPE, parmi, parmj, n) \
    {                                   \
        size_t i = (n) / sizeof(TYPE);  \
        TYPE *pi = (TYPE *)(parmi);     \
        TYPE *pj = (TYPE *)(parmj);     \
        do                              \
        {                               \
            TYPE t = *pi;               \
            *pi++ = *pj;                \
            *pj++ = t;                  \
        } while (--i > 0);              \
    }

static inline void
swapfunc(char *a, char *b, size_t n, int swaptype_long, int swaptype_int)
{
    if (swaptype_long <= 1)
        swapcode(long, a, b, n) else if (swaptype_int <= 1)
            swapcode(int, a, b, n) else swapcode(char, a, b, n)
}

#define swap(a, b)                   \
    if (swaptype_long == 0)          \
    {                                \
        long t = *(long *)(a);       \
        *(long *)(a) = *(long *)(b); \
        *(long *)(b) = t;            \
    }                                \
    else if (swaptype_int == 0)      \
    {                                \
        int t = *(int *)(a);         \
        *(int *)(a) = *(int *)(b);   \
        *(int *)(b) = t;             \
    }                                \
    else                             \
        swapfunc(a, b, es, swaptype_long, swaptype_int)

#define vecswap(a, b, n) \
    if ((n) > 0)         \
    swapfunc(a, b, n, swaptype_long, swaptype_int)

#define CMP(t, x, y) (cmp((x), (y)))

/*
 * Find the median of 3 elements
 */
static inline char *
med3(char *a, char *b, char *c, cmp_t *cmp, void *thunk)
{
    return CMP(thunk, a, b) < 0 ? (CMP(thunk, b, c) < 0 ? b : (CMP(thunk, a, c) < 0 ? c : a))
                                : (CMP(thunk, b, c) > 0 ? b : (CMP(thunk, a, c) < 0 ? a : c));
}

#define DEPTH(x) (2 * (fls((int)(x)) - 1))

struct thin_file
{
    char *name;
    cpu_type_t cputype;
    cpu_subtype_t cpusubtype;
    uint32_t align;
};

void dump(void *a, size_t n, size_t es)
{
    printf("dump-----\n");
    for (int i = 0; i < n; i++)
    {
        struct thin_file *c = a + i * es;
        printf("%s\n", c->name);
    }
    printf("dump-----\n");
}

/*
 * Simple insertion sort routine.
 */
static bool
_isort(void *a, size_t n, size_t es, void *thunk, cmp_t *cmp, int swap_limit, int swaptype_long, int swaptype_int)
{
    int swap_cnt = 0;
    for (char *pm = (char *)a + es; pm < (char *)a + n * es; pm += es)
    {
        for (char *pl = pm; pl > (char *)a && CMP(thunk, pl - es, pl) > 0;
             pl -= es)
        {
            swap(pl, pl - es);
            if (swap_limit && ++swap_cnt > swap_limit)
                return false;
        }
    }
    return true;
}

#define thunk NULL
static void
_qsort(void *a, size_t n, size_t es, cmp_t *cmp, int depth_limit)

{
    char *pa, *pb, *pc, *pd, *pl, *pm, *pn;
    size_t d1, d2;
    int cmp_result;
    int swaptype_long, swaptype_int, swap_cnt;
    printf("myqsort is called %d\n", depth_limit);
loop:
    // SWAPINIT(long, a, es);
    // SWAPINIT(int, a, es);
    swap_cnt = 0;
    printf("num %d\n", n);
    if (depth_limit-- <= 0)
    {
        /*
         * We've hit our recursion limit, switch to heapsort
         */
        printf("switch to myheapsort\n");
        myheapsort(a, n, es, cmp);

        return;
    }

    if (n <= 7)
    {
        /*
         * For sufficiently small inputs, we'll just insertion sort.
         *
         * Pass 0 as swap limit, since this must complete.
         */
        printf("switch to isort\n");
        _isort(a, n, es, thunk, cmp, 0, swaptype_long, swaptype_int);
        dump(a, n, es);
        return;
    }

    /*
     * Compute the pseudomedian.  Small arrays use 3 samples, large ones use 9.
     */
    pl = a;
    pm = (char *)a + (n / 2) * es;
    pn = (char *)a + (n - 1) * es;
    if (n > 40)
    {
        size_t d = (n / 8) * es;

        pl = med3(pl, pl + d, pl + 2 * d, cmp, thunk);
        pm = med3(pm - d, pm, pm + d, cmp, thunk);
        pn = med3(pn - 2 * d, pn - d, pn, cmp, thunk);
    }
    struct thin_file *t;
    t = pl;
    printf("pl %s\n", t->name);
    t = pm;
    printf("pm %s\n", t->name);
    t = pn;
    printf("pn %s\n", t->name);

    pm = med3(pl, pm, pn, cmp, thunk);
    t = pm;
    printf("med3 pm %s\n", t->name);

    /*
     * Pull the median to the front, starting us with:
     *
     * +-+-------------+
     * |=|      ?      |
     * +-+-------------+
     * a pa,pb         pc,pd
     */
    swap(a, pm);
    pa = pb = (char *)a + es;
    pc = pd = (char *)a + (n - 1) * es;

    for (;;)
    {
        /*
         * - Move b forward while it's less than the median
         * - Move c backwards while it's greater than the median
         * - When equal to the median, swap to the outside
         */
        while (pb <= pc && (cmp_result = CMP(thunk, pb, a)) <= 0)
        {
            if (cmp_result == 0)
            {
                printf("pb <= pc cmp_result=0\n");
                swap_cnt = 1;
                swap(pa, pb);
                pa += es;
            }
            printf("pb <= pc bp++\n");
            pb += es;
        }
        while (pb <= pc && (cmp_result = CMP(thunk, pc, a)) >= 0)
        {
            if (cmp_result == 0)
            {
                printf("pb <= pc cmp_result=0\n");
                swap_cnt = 1;
                swap(pc, pd);
                pd -= es;
            }
            printf("pb <= pc pc--\n");
            pc -= es;
        }
        if (pb > pc)
            break;
        swap(pb, pc);
        swap_cnt = 1;
        pb += es;
        pc -= es;
    }

    /*
     * Now we've got:
     *
     * +---+-----+-----+---+
     * | = |  <  |  >  | = |
     * +---+-----+-----+---+
     * a   pa  pc,pb   pd  pn
     *
     * So swap the '=' into the middle
     */

    pn = (char *)a + n * es;
    d1 = MIN(pa - (char *)a, pb - pa);
    vecswap(a, pb - d1, d1);
    d1 = MIN(pd - pc, pn - pd - es);
    vecswap(pb, pn - d1, d1);

    /*
     * +-----+---+---+-----+
     * |  <  |   =   |  >  |
     * +-----+---+---+-----+
     * a                   pn
     */

    if (swap_cnt == 0)
    {                      /* Switch to insertion sort */
        int r = 1 + n / 4; /* n > 7, so r >= 2 */
        if (!_isort(a, n, es, thunk, cmp, r, swaptype_long, swaptype_int))
        {
            printf("goto nevermind;\n");
            goto nevermind;
        }
        printf("return swap_cnt=0 isort");
        return;
    }
nevermind:

    d1 = pb - pa;
    d2 = pd - pc;
    if (d1 <= d2)
    {
        /* Recurse on left partition, then iterate on right partition */
        if (d1 > es)
        {
            printf("d1 > 1 do qsort\n");
            _qsort(a, d1 / es, es, cmp, depth_limit);
        }
        if (d2 > es)
        {
            printf("d2 > 1 goto loop\n");
            /* Iterate rather than recurse to save stack space */
            /* qsort(pn - d2, d2 / es, es, cmp); */
            a = pn - d2;
            n = d2 / es;
            goto loop;
        }
    }
    else
    {
        /* Recurse on right partition, then iterate on left partition */
        if (d2 > es)
        {
            printf("d2 > 1 do qsort\n");
            _qsort(pn - d2, d2 / es, es, cmp, depth_limit);
        }
        if (d1 > es)
        {
            printf("d1 > 1 goto loop\n");
            /* Iterate rather than recurse to save stack space */
            /* qsort(a, d1 / es, es, cmp); */
            n = d1 / es;
            goto loop;
        }
    }
}

void qsort(void *a, size_t n, size_t es, cmp_t *cmp)
{
    _qsort(a, n, es,
           cmp, DEPTH(n));
}
