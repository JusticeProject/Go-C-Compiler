package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

/////////////////////////////////////////////////////////////////////////////////

func main() {
	if runtime.GOOS == "linux" {
		doLinux()
	} else {
		doWindows()
	}
}

/////////////////////////////////////////////////////////////////////////////////

func doLinux() {
	fmt.Println("running Linux version")
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: goc /path/to/source.c")
		fmt.Println("Options:")
		fmt.Println("--lex will run lexer but stop before parsing, no output files are produced")
		fmt.Println("--parse will run lexer and parser but stop before assembly generation, no output files are produced")
		fmt.Println("--codegen will run lexer, parser, and assembly generation but stop before code emission, no output files are produced")
		fmt.Println("-S will emit an assembly file but will not assemble or link it")
		os.Exit(1)
	}

	inputFilename := os.Args[1]
	if filepath.Ext(inputFilename) != ".c" {
		fmt.Println("Please use a file with .c extension")
		os.Exit(1)
	}
	fmt.Println("found inputFilename", inputFilename)

	runParser := true
	runAssemblyGeneration := true
	runCodeEmission := true
	produceExecutable := true

	fmt.Println("found", len(os.Args), "args")

	if len(os.Args) == 3 {
		switch os.Args[2] {
		case "--lex":
			fmt.Println("stopping after lexer")
			runParser = false
			runAssemblyGeneration = false
			runCodeEmission = false
			produceExecutable = false
		case "--parse":
			fmt.Println("stopping after parser")
			runParser = true
			runAssemblyGeneration = false
			runCodeEmission = false
			produceExecutable = false
		case "--codegen":
			fmt.Println("stopping after assembly generation")
			runParser = true
			runAssemblyGeneration = true
			runCodeEmission = false
			produceExecutable = false
		case "-S":
			fmt.Println("stopping after code emission")
			runParser = true
			runAssemblyGeneration = true
			runCodeEmission = true
			produceExecutable = false
		default:
			fmt.Println("unknown option, exiting")
			os.Exit(1)
		}
	}

	// produce preprocessed file
	preprocessedFilename := strings.TrimSuffix(inputFilename, ".c") + ".i"
	outBytes, err := exec.Command("gcc", "-E", "-P", inputFilename, "-o", preprocessedFilename).CombinedOutput()
	if err != nil {
		fmt.Println("gcc returned error:", err)
		fmt.Printf("additional info: %s\n", outBytes)
		os.Exit(1)
	}

	fileContents := loadFile(preprocessedFilename)

	// delete preprocessed file
	err = os.Remove(preprocessedFilename)
	if err != nil {
		fmt.Println("Could not remove preprocessedFile")
	}

	// do the compilation
	assemblyFilename := strings.TrimSuffix(inputFilename, ".c") + ".s"
	doFourCompilerSteps(fileContents, runParser, runAssemblyGeneration, runCodeEmission, assemblyFilename)

	if !produceExecutable {
		os.Exit(0)
	}

	// assembly and link using gcc (produce executable)
	binaryFilename := strings.TrimSuffix(inputFilename, ".c")
	outBytes, err = exec.Command("gcc", assemblyFilename, "-o", binaryFilename).CombinedOutput()
	if err != nil {
		fmt.Println("gcc returned error:", err)
		fmt.Printf("additional info: %s\n", outBytes)
		os.Exit(1)
	}

	// remove the assembly file
	err = os.Remove(assemblyFilename)
	if err != nil {
		fmt.Println("Could not remove assemblyFilename")
	}
}

/////////////////////////////////////////////////////////////////////////////////

func doWindows() {
	fmt.Println("running Windows debug version")
	filename := "return_2.c"

	contents := loadFile(filename)
	assemblyFilename := strings.TrimSuffix(filename, ".c") + ".s"
	doFourCompilerSteps(contents, true, true, true, assemblyFilename)
}

/////////////////////////////////////////////////////////////////////////////////

func doFourCompilerSteps(fileContents string, runParser bool, runAssemblyGeneration bool, runCodeEmission bool, assemblyFilename string) {
	fmt.Println("running compiler with fileContents:")
	fmt.Println(fileContents)

	// run lexer
	tokens := doLexer(fileContents)
	fmt.Println("found tokens:")
	fmt.Println(tokens)

	if !runParser {
		fmt.Println("not running parser, done")
		os.Exit(0)
	}

	// run parser, get the Abstract Syntax Tree
	ast := doParser(tokens)

	if !runAssemblyGeneration {
		os.Exit(0)
	}

	// run assembly generation
	asm := doAssemblyGen(ast)

	if !runCodeEmission {
		os.Exit(0)
	}

	//run code emission
	doCodeEmission(asm, assemblyFilename)
}

/////////////////////////////////////////////////////////////////////////////////

func loadFile(filename string) string {
	data, _ := os.ReadFile(filename)
	return string(data)
}
