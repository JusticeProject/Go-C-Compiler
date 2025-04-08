package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

/////////////////////////////////////////////////////////////////////////////////

func main() {
	fmt.Println("Using Go version", runtime.Version())

	if runtime.GOOS == "linux" {
		doLinux()
	} else {
		doWindows()
	}
}

/////////////////////////////////////////////////////////////////////////////////

func doLinux() {
	fmt.Println("running Linux version")
	if len(os.Args) < 2 {
		fmt.Println("Usage: goc /path/to/source.c")
		fmt.Println("Options:")
		fmt.Println("--lex will stop after the lexer and before the parser, no output files are produced")
		fmt.Println("--parse will stop after the parser and before semantic analysis, no output files are produced")
		fmt.Println("--validate will stop after semantic analysis and before tacky generation, no output files are produced")
		fmt.Println("--tacky will stop after tacky generation and before assembly generation, no output files are produced")
		fmt.Println("--codegen will stop after assembly generation and before code emission, no output files are produced")
		fmt.Println("-S will emit an assembly file but will not assemble or link it")
		fmt.Println("-c will emit an object file but will not link it")
		fmt.Println("-o is used to specify the executable name. The default is to use the first .c file and remove the .c from the name.")
		os.Exit(1)
	}

	fmt.Println("found", len(os.Args), "args")

	allInputFileNames := []string{}
	libraries := []string{}

	runParser := true
	runSemanticAnalysis := true
	runTackyGeneration := true
	runAssemblyGeneration := true
	runCodeEmission := true

	produceObjectFile := false
	produceExecutable := true
	outputFilename := ""

	// index 0 is the program currently running (./goc)
	for index := 1; index < len(os.Args); index++ {
		currentArg := os.Args[index]

		if filepath.Ext(currentArg) == ".c" {
			allInputFileNames = append(allInputFileNames, currentArg)
			fmt.Println("found source file", currentArg)
		}

		if currentArg[0:1] == "-" {
			switch currentArg {
			case "--lex":
				fmt.Println("stopping after lexer")
				runParser = false
			case "--parse":
				fmt.Println("stopping after parser")
				runSemanticAnalysis = false
			case "--validate":
				fmt.Println("stopping after semantic analysis")
				runTackyGeneration = false
			case "--tacky":
				fmt.Println("stopping after tacky generation")
				runAssemblyGeneration = false
			case "--codegen":
				fmt.Println("stopping after assembly generation")
				runCodeEmission = false
			case "-S":
				fmt.Println("stopping after code emission")
				produceExecutable = false
			case "-c":
				fmt.Println("creating object file instead of executable")
				produceObjectFile = true
				produceExecutable = false
			case "-o":
				// check that index + 1 is valid before using it
				if (index + 1) < len(os.Args) {
					outputFilename = os.Args[index+1]
					index++
				}
			default:
				// it could be a library that we need to link
				re, _ := regexp.Compile(`-l[a-zA-Z0-9]+`)
				result := re.FindStringIndex(currentArg)
				if len(result) > 0 {
					lib := currentArg[result[0]:result[1]]
					fmt.Println("Will link library using", lib)
					libraries = append(libraries, lib)
				} else {
					fail("unknown option, exiting")
				}
			}
		}
	}

	allAssemblyFilenames := []string{}
	for _, filename := range allInputFileNames {
		// produce preprocessed file
		preprocessedFilename := strings.TrimSuffix(filename, ".c") + ".i"
		outBytes, err := exec.Command("gcc", "-E", "-P", filename, "-o", preprocessedFilename).CombinedOutput()
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

		// do the compilation and produce an assembly file
		assemblyFilename := strings.TrimSuffix(filename, ".c") + ".s"
		allAssemblyFilenames = append(allAssemblyFilenames, assemblyFilename)
		doCompilerSteps(fileContents, runParser, runSemanticAnalysis, runTackyGeneration, runAssemblyGeneration, runCodeEmission, assemblyFilename)
	}

	if produceObjectFile {
		// assemble but don't link for each assembly file
		for _, assemblyFile := range allAssemblyFilenames {
			objectFilename := strings.TrimSuffix(assemblyFile, ".s") + ".o"
			outBytes, err := exec.Command("gcc", "-c", assemblyFile, "-o", objectFilename).CombinedOutput()
			if err != nil {
				fmt.Println("gcc returned error:", err)
				fmt.Printf("additional info: %s\n", outBytes)
				os.Exit(1)
			}

			fmt.Println("object file created:", objectFilename)
		}
	} else if produceExecutable {
		// assembly and link using gcc
		fmt.Println("running assembler and linker")
		gccArgs := make([]string, len(allAssemblyFilenames))
		copy(gccArgs, allAssemblyFilenames)
		gccArgs = append(gccArgs, "-o")
		if outputFilename == "" {
			// no output filename was given, so use the first .c file
			outputFilename = strings.TrimSuffix(allInputFileNames[0], ".c")
		}
		gccArgs = append(gccArgs, outputFilename)
		gccArgs = append(gccArgs, libraries...)
		fmt.Println("Running cmd gcc with args:", gccArgs)

		outBytes, err := exec.Command("gcc", gccArgs...).CombinedOutput()
		if err != nil {
			fmt.Println("gcc returned error:", err)
			fmt.Printf("additional info: %s\n", outBytes)
			os.Exit(1)
		}
		fmt.Printf("additional info: %s\n", outBytes)
		fmt.Println("executable created:", outputFilename)
	}

	// remove the assembly file(s)
	for _, filename := range allAssemblyFilenames {
		err := os.Remove(filename)
		if err != nil {
			fmt.Println("Could not remove assembly file", filename)
		}
	}

}

/////////////////////////////////////////////////////////////////////////////////

func doWindows() {
	fmt.Println("running Windows debug version")
	filename := "test.c"

	contents := loadFile(filename)
	assemblyFilename := strings.TrimSuffix(filename, ".c") + ".s"
	doCompilerSteps(contents, true, true, true, true, true, assemblyFilename)
}

/////////////////////////////////////////////////////////////////////////////////

func doCompilerSteps(fileContents string, runParser bool, runSemanticAnalysis bool, runTackyGeneration bool,
	runAssemblyGeneration bool, runCodeEmission bool, assemblyFilename string) {

	fmt.Println("running compiler with fileContents:")
	fmt.Println(fileContents)

	// run lexer
	fmt.Println("running lexer")
	tokens := doLexer(fileContents)
	fmt.Println("found tokens:")
	fmt.Println(tokens)

	if !runParser {
		fmt.Println("not running parser, done")
		os.Exit(0)
	}

	// run parser, get the Abstract Syntax Tree
	fmt.Println("running parser")
	ast := doParser(tokens)

	if !runSemanticAnalysis {
		fmt.Println("not running semantic analysis, done")
		os.Exit(0)
	}

	// run semantic analysis and update the Abstract Syntax Tree
	fmt.Println("running semantic analysis")
	ast = doIdentifierResolution(ast)
	ast = doTypeChecking(ast)
	ast = doLoopLabeling(ast)

	if !runTackyGeneration {
		fmt.Println("not running tacky generation, done")
		os.Exit(0)
	}

	// run tacky generation
	fmt.Println("running tacky generation")
	tacky := doTackyGen(ast)

	if !runAssemblyGeneration {
		os.Exit(0)
	}

	// run assembly generation
	fmt.Println("running assembly generation")
	asm := doAssemblyGen(tacky)

	if !runCodeEmission {
		os.Exit(0)
	}

	//run code emission
	fmt.Println("running code emission")
	doCodeEmission(asm, assemblyFilename)
}

/////////////////////////////////////////////////////////////////////////////////

func loadFile(filename string) string {
	data, _ := os.ReadFile(filename)
	return string(data)
}

/////////////////////////////////////////////////////////////////////////////////

func fail(msg ...string) {
	joinedMsg := strings.Join(msg, " ")
	fmt.Println(joinedMsg)
	os.Exit(1)
}
