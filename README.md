# Go-C-Compiler
A C compiler implemented in Go. Runs on Linux. Based on the book [Writing a C Compiler](https://nostarch.com/writing-c-compiler) by Nora Sandler. Currently only supports the main features presented in the book (not the extra credit features).

# Usage
First, install the go compiler.
```bash
sudo apt install gccgo
```
Second, compile the C compiler.
```bash
cd chapter9
gccgo *.go -o goc
```
Third, use the new C compiler to compile a .c file.
```
./goc test.c
```
Fourth, run the executable.
```bash
./test
```
