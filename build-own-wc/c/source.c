#include <stdio.h>
#include <wchar.h>
#include <stdlib.h>
#include <string.h>
#include <locale.h>
#include <unistd.h>

int count_words(char* value) {
    int result = 0;
    int inword = 0;

    do switch (*value) {
        case '\0': case ' ': case '\t': case '\n': case '\r':
            if (inword) {
                inword = 0;
                result++;
            }
            break;
        default:
            inword = 1;
            break;
    } while (*value++);

    return result;
}

int main(int argc, char *argv[]) {
    char *filename = NULL;
    char *command = NULL;
    FILE *fptr = NULL;

    // Set locale for multi-byte char support
    setlocale(LC_ALL, "");

    if (argc > 3) {
        printf("Argument cannot more than one\n");
        return 1;
    } else if (argc < 2) {
        printf("Argument is not valid, please enter argument\n");
        return 1;
    }

    if (strcmp(argv[1], "-c") == 0) {
        command = argv[1];
    } else if (strcmp(argv[1], "-l") == 0) {
        command = argv[1];
    } else if (strcmp(argv[1], "-w") == 0) {
        command = argv[1];
    } else if (strcmp(argv[1], "-m") == 0) {
        command = argv[1];
    } else {
        filename = argv[1];
    }

    if (argc == 3) {
        filename = argv[2];
    }

    if (filename != NULL) {
        fptr = fopen(filename, "r");
    } else if (!isatty(fileno(stdin))) {
        fptr = stdin;
    }

    if (fptr == NULL) {
        printf("Cannot find the file, please enter the valid path.\n");
        return 1;
    }

    int result = 0;
    int lines, characters, words = 0;
    char content[256];

    if (command != NULL) {
        if (strcmp(command, "-m") == 0) {
            while (fgetwc(fptr) != WEOF) {
                result++;
            }
        } else {
            while (fgets(content, sizeof(content), fptr) != NULL) {
                if (strcmp(command, "-l") == 0) {
                    result++;
                } else if (strcmp(command, "-c") == 0) {
                    result += strlen(content);
                } else if (strcmp(command, "-w") == 0) {
                    result += count_words(content);
                }
            }
        }
    } else {
        while (fgets(content, sizeof(content), fptr) != NULL) {
            lines++;
            characters += strlen(content);
            words += count_words(content);
        }
    }

    if (filename == NULL) {
        printf("  %d\n", result);
    } else if (command != NULL) {
        printf("  %d %s\n", result, filename);
    } else {
        printf("  %d %d %d %s\n", lines, words, characters, filename);
    }

    fclose(fptr);

    return 0;
}
