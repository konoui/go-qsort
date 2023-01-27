// +build ignore

#include <mach/mach.h>
#include <stdio.h>
#include <errno.h>

int myheapsort(void *__base, size_t __nel, size_t __width,
               int (*_Nonnull __compar)(const void *, const void *));

struct thin_file
{
    char *name;
    cpu_type_t cputype;
    cpu_subtype_t cpusubtype;
    uint32_t align;
};

static int
cmp_qsort(
    const struct thin_file *thin1,
    const struct thin_file *thin2)
{
    /* if cpu types match, sort by cpu subtype */
    if (thin1->cputype == thin2->cputype)
        return ((thin1->cpusubtype & ~CPU_SUBTYPE_MASK) -
                (thin2->cpusubtype & ~CPU_SUBTYPE_MASK));

    /* force arm64-family to follow after all other slices */
    if (thin1->cputype == CPU_TYPE_ARM64)
        return 1;
    if (thin2->cputype == CPU_TYPE_ARM64)
        return -1;

    /* sort all other cpu types by alignment */
    return thin1->align - thin2->align;
}

int main()
{
    struct thin_file files[] =
        {
            {.name = "i386", .cputype = 7, .cpusubtype = 3, .align = 0},
            {.name = "x86_64", .cputype = 16777223, .cpusubtype = 3, .align = 0},
            {.name = "x86_64h", .cputype = 16777223, .cpusubtype = 8, .align = 0},
            {.name = "arm", .cputype = 12, .cpusubtype = 0, .align = 0},
            {.name = "armv4t", .cputype = 12, .cpusubtype = 5, .align = 0},
            {.name = "armv6", .cputype = 12, .cpusubtype = 6, .align = 0},
            {.name = "armv7", .cputype = 12, .cpusubtype = 9, .align = 0},
            {.name = "armv7f", .cputype = 12, .cpusubtype = 10, .align = 0},
            {.name = "armv7s", .cputype = 12, .cpusubtype = 11, .align = 0},
            {.name = "armv7k", .cputype = 12, .cpusubtype = 12, .align = 0},
            {.name = "armv6m", .cputype = 12, .cpusubtype = 14, .align = 0},
            {.name = "armv7m", .cputype = 12, .cpusubtype = 15, .align = 0},
            {.name = "armv7em", .cputype = 12, .cpusubtype = 16, .align = 0},
            {.name = "armv8m", .cputype = 12, .cpusubtype = 17, .align = 0},
            {.name = "arm64", .cputype = 16777228, .cpusubtype = 0, .align = 0},
            {.name = "arm64e", .cputype = 16777228, .cpusubtype = 2, .align = 0},
            {.name = "arm64v8", .cputype = 16777228, .cpusubtype = 1, .align = 0},
            {.name = "arm64_32", .cputype = 33554444, .cpusubtype = 0, .align = 0},
        };
    size_t array_size = sizeof(files) / sizeof(files[0]);

    struct thin_file qfiles[array_size];
    memcpy(qfiles, files, sizeof(files));

    // printf("heapsort----------------------------------------------------------------------\n");
    // int ret = myheapsort(&files, array_size, sizeof(files[0]), cmp_qsort);
    // if (ret != 0)
    // {
    //     fprintf(stderr, "heapsort abort ret is non zero: %d %d\n", ret, errno);
    //     return ret;
    // }
    // for (int i = 0; i < array_size; i++)
    //     printf("%s\n", files[i].name);

    printf("qsort----------------------------------------------------------------------\n");
    qsort(&qfiles, array_size, sizeof(qfiles[0]), cmp_qsort);
    for (int i = 0; i < array_size; i++)
        printf("%s\n", qfiles[i].name);
}
