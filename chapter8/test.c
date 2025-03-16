int main(void) {
    int sum = 0;

    for (int i = 1; i <=4; i = i + 1)
    {
        sum = sum + i;
    }

    int a = 1;
    while (1) {
        sum = sum + a;
        a = a + 1;

        if (a > 4) {
            break;
        }
    }

    a = 0;
    do {
        a = a + 1;
        if (a > 4) {
            continue;
        }
        sum = sum + a;
    } while (a <= 4);

    return sum;
}
