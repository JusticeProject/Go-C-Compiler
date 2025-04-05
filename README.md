# Go-C-Compiler
A C compiler implemented in Go. Runs on Linux. Based on the book [Writing a C Compiler](https://nostarch.com/writing-c-compiler) by Nora Sandler. Currently only supports the main features presented in the book (not the extra credit features). The commit comments will indicate if a particular chapter has been completed or not.

# Usage
First, [install the go compiler on Linux](https://go.dev/doc/install).

Second, compile the C compiler. Chapter 9 has been completed successfully with all tests passing so we'll use that code.
```bash
cd chapter9
go build -o goc *.go
```
Third, use the new C compiler to compile a .c file.
```
./goc test.c
```
Fourth, run the executable.
```bash
./test
```
