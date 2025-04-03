int scale(int input, double scaleFactor);

int main(void) {
    result = scale(6, 1.51);
    return result;
}

int scale(int input, double scaleFactor) {
    return input * scaleFactor;
}