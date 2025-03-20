# Go-C-Compiler
A C compiler implemented in Go. Based on the book [Writing a C Compiler](https://nostarch.com/writing-c-compiler) by Nora Sandler.

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
