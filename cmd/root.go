package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	p "path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "wb",
	Short: "read whiteboard",
	Run:   rootRun,
}

func init() {
}

var (
	out      string //holds coggle data as a string
	rootDir  string //path of the folder containing the working path
	rootName string //current folder name
	curPath  string //current path
)

func rootRun(cmd *cobra.Command, args []string) {

	var err error
	curPath, err = filepath.Abs("")
	if err != nil {
		fmt.Println(err.Error())
	}

	rootName = p.Base(curPath)
	coggleFileName := rootName + ".txt"

	rootDir = p.Dir(curPath)
	//loop through all the files
	err = filepath.Walk(curPath, visit)
	if err != nil {
		fmt.Println(err.Error())
	}

	//write final output to file
	//fmt.Println(out)
	err = ioutil.WriteFile(coggleFileName, []byte(out), 0644)
	if err != nil {
		fmt.Println(err.Error())
	}

	//cmd output
	txtOutPath := p.Join(curPath, coggleFileName)
	fmt.Println("New coggle file written to : \n" + txtOutPath)
}

func visit(path string, f os.FileInfo, err error) error {
	if strings.Contains(path, ".git") {
		return nil
	}
	if strings.Contains(path, "vendor") {
		return nil
	}

	working := strings.Replace(path, rootDir, "", 1)

	//return if not a path or .go file
	filePath := strings.Replace(path, curPath+"/", "", 1)
	var workingIsFile bool = false
	if len(filePath) > 0 {
		var err error
		workingIsFile, err = isFile(filePath)
		if err != nil {
			return err
		}

		if workingIsFile && !strings.Contains(filePath, ".go") {
			return nil
		}
	}

	finalElement := p.Base(path) + "\n"

	//if file loop through the contents and add them to the file
	if workingIsFile {

		//get the file contents
		fileContentsByte, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		fileContents := string(fileContentsByte)

		//parse the file
		fset := token.NewFileSet() // positions are relative to fset

		// Parse the file string
		f, err := parser.ParseFile(fset, "", fileContents, 0) // parser.ImportsOnly)
		if err != nil {
			fmt.Println(err)
		}

		for i := 0; i < len(f.Decls); i++ {

			switch f.Decls[i].(type) {
			case *ast.GenDecl:
				typeDecl := f.Decls[i].(*ast.GenDecl)
				for k := 0; k < len(typeDecl.Specs); k++ {
					switch typeDecl.Specs[k].(type) {
					case *ast.TypeSpec:
						writeTypeOutput := false
						typeSpec := typeDecl.Specs[k].(*ast.TypeSpec)
						typeType := typeSpec.Type

						var start, end token.Pos

						switch typeType.(type) {
						case *ast.StructType:
							structDecl := typeType.(*ast.StructType)
							fields := structDecl.Fields.List

							//define start position
							start = structDecl.Struct - 1

							// Define the end position of final field
							for _, field := range fields {
								end = field.Type.End() - 1

							}
							writeTypeOutput = true

						case *ast.InterfaceType:
							structDecl := typeType.(*ast.InterfaceType)
							fields := structDecl.Methods.List

							//define start position
							start = structDecl.Interface - 1

							// Define the end position of final field
							for _, field := range fields {
								end = field.Type.End() - 1

							}
							writeTypeOutput = true

						}
						if writeTypeOutput {
							toAdd := typeSpec.Name.Name + " " +
								fileContents[start:end] +
								"}"
							toAdd = strings.Replace(toAdd, "\r\n", ", ", -1)
							toAdd = strings.Replace(toAdd, "\n", ", ", -1)
							toAdd = strings.Replace(toAdd, "\t", "", -1)
							toAdd = strings.Replace(toAdd, "{,", "{ ", -1)
							for j := 0; j < 10; j++ {
								toAdd = strings.Replace(toAdd, "  ", " ", -1)
							}
							finalElement += "\t" + toAdd + "\n"
						}
					}
				}
			case *ast.FuncDecl:
				typeDecl := f.Decls[i].(*ast.FuncDecl)
				finalElement += "\tfunc " + typeDecl.Name.Name + "\n"
			}
		}
	}

	//generate appropriate tabs
	numTabs := strings.Count(working, "/") - 1

	var tabs string
	for i := 0; i < numTabs; i++ {
		tabs += "\t"
	}

	scanner := bufio.NewScanner(strings.NewReader(finalElement))
	for scanner.Scan() {
		out += tabs + scanner.Text() + "\n"
	}

	return nil
}

func isFile(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return true, err
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return false, nil
	case mode.IsRegular():
		return true, nil
	}

	return true, errors.New("not file or dir")
}
