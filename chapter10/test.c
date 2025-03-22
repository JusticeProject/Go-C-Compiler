int putchar(int c);

int fibonacci(void) {
    static int n = -1;
    static int np1 = 1;

    int np2 = n + np1;

    n = np1;
    np1 = np2;

    return np2;
}

int main(void) {
    int result = 0;
    for (int i = 0; i <= 6; i=i+1) {
        result = fibonacci();
        putchar(result + 48);
    }
    putchar(10);
    return result;
}
