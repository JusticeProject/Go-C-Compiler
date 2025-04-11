int *ptr = 0;

int main(void) {
    int x = 0;
    ptr = &x;
    *ptr = 4;
    return *ptr;
}
