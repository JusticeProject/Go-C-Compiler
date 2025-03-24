int putchar(int c);

int n;
int np1;

int fibonacci(void) {
    int np2 = n + np1;

    n = np1;
    np1 = np2;

    return np2;
}

int main(void) {
    n = -1;
    np1 = 1;
    int result = 0;

    for (int i = 0; i <= 6; i=i+1) {
        result = fibonacci();
        putchar(result + 48);
    }
    putchar(10);
    return result;
}
