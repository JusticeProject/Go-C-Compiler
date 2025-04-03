static int test = 2.8;

int scale(int input, double scaleFactor);

int main(void) {
    int result = scale(6, 1.51);
    return result + test;
}

int scale(int input, double scaleFactor) {
    return input * scaleFactor;
}