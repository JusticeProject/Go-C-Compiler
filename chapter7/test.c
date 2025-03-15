int main(void) {
    int a = 0;

    if (a < 1)
        a = a + 1;

    if (a < 1)
        a = 1;
    else
        a = 10;

    a = (a == 10) ? 20 : 30;
    a = (a == 19) ? 40 : 60;

    return a;
}
