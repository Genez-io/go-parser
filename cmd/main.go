// extract_comments.go
package main

import (
	"encoding/json"
	"fmt"
	"gnz-go-parser/models"
	"go/ast"
	"go/parser"
	"go/token"

	"os"
	"strings"
)
type Error struct {
	Error string `json:"error"`
}

type Response struct {
	Classes []models.Class `json:"classes"`
	Methods []models.Method `json:"methods"`
}

func SendError(err error) {
	json, err := json.Marshal(Error{
		Error: err.Error(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json))
}

// Function to extract comments in front of type and method declarations
func extractComments(filePath string) (*Response, error) {
    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
    if err != nil {
        return nil, err
    }

	var classes []models.Class 
	var methods []models.Method

    // Inspect the AST and collect comments
    ast.Inspect(node, func(n ast.Node) bool {
        switch x := n.(type) {
        case *ast.GenDecl:
            if x.Tok == token.TYPE && x.Doc != nil {
                if len(x.Doc.List) >0 {
					lastComment := x.Doc.List[len(x.Doc.List)-1]
					if strings.Contains(lastComment.Text, "genezio: deploy") {
						var class models.Class
						class.Comment = lastComment.Text
						class.Name = x.Specs[0].(*ast.TypeSpec).Name.Name
						classes = append(classes, class)
					}
				}
            }
        case *ast.FuncDecl:
			if x.Recv != nil && x.Doc != nil {
				// Extract the receiver type (class name)
				var className string
				for _, field := range x.Recv.List {
					// Dereference the pointer type, if necessary
					switch expr := field.Type.(type) {
					case *ast.Ident:
						className = expr.Name
					case *ast.StarExpr:
						if ident, ok := expr.X.(*ast.Ident); ok {
							className = ident.Name
						}
					}
				}

				if className != "" {
					 if len(x.Doc.List) > 0 {
						lastComment := x.Doc.List[len(x.Doc.List)-1]
						if strings.Contains(lastComment.Text, "genezio:") {
							var method models.Method
							method.Comment = lastComment.Text
							method.Name = x.Name.Name
							method.ClassName = className
							methods = append(methods, method)
						}
					 }
				}

			}
        }
        return true
    })

	return &Response{
		Classes: classes,
		Methods: methods,
	},nil
}

func main() {

    filePath := os.Args[1] // Replace with your file path

		
    response, err := extractComments(filePath)
    if err != nil {
        SendError(err)
		return
    }

	json, err := json.MarshalIndent(response,"","   ")
	if err != nil {
		SendError(err)
		return
	}
	fmt.Println(string(json))

}