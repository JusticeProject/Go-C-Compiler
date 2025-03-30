unsigned long sum(unsigned long a, unsigned long b);

int main(void) {
    unsigned long result = sum(9UL, (unsigned long)2);
    return result;
}

unsigned long sum(unsigned long a, unsigned long b) {
    return a + b;
}