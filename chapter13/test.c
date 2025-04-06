static unsigned int test = 2.8;

unsigned int scale(unsigned int input, double scaleFactor);

int main(void) {
    unsigned int result = scale(6, 1.51);
    return result + test;
}

unsigned int scale(unsigned int input, double scaleFactor) {
    return input * scaleFactor;
}