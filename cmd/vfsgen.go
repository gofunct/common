package cmd

import (
	"github.com/spf13/cobra"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

var (
	outFile    string
	outPackage string
	tags       string
	varName    string
	varComment string
	fs         http.FileSystem = http.Dir(tmplDir)
	tmplDir    string
)

// Initialize
func init() {
	{ // Commands

	}

	{ // Flags
		vfsGenCommand.Flags().StringVarP(&outFile, "out-f", "f", "vfsgen.go", "output file name")
		vfsGenCommand.Flags().StringVarP(&outPackage, "out-p", "p", "templates", "output package name")
		vfsGenCommand.Flags().StringVarP(&tags, "tags", "t", "vfsgen", "build tags")
		vfsGenCommand.Flags().StringVarP(&varName, "var", "v", "FS", "variable name")
		vfsGenCommand.Flags().StringVarP(&varComment, "comment", "c", "", "variable comment")
		vfsGenCommand.Flags().StringVarP(&tmplDir, "tmpl-dir", "d", "_data", "template directory")
	}
}

var vfsGenCommand = &cobra.Command{
	Use:     "vfsgen",
	Aliases: []string{"gen", "template"},
	Short:   "vfsgen specified templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vfsgen.Generate(fs, vfsgen.Options{

			// Filename of the generated Go code output (including extension).
			// If left empty, it defaults to "{{toLower .VariableName}}_vfsdata.go".
			Filename: outFile,

			// PackageName is the name of the package in the generated code.
			// If left empty, it defaults to "main".
			PackageName: outPackage,

			// BuildTags are the optional build tags in the generated code.
			// The build tags syntax is specified by the go tool.
			BuildTags: tags,
			// VariableName is the name of the http.FileSystem variable in the generated code.
			// If left empty, it defaults to "assets".
			VariableName: varName,
			// VariableComment is the comment of the http.FileSystem variable in the generated code.
			// If left empty, it defaults to "{{.VariableName}} statically implements the virtual filesystem provided to vfsgen.".
			VariableComment: varComment,
		})
		if err != nil {
			return err
		}
		return nil
	},
}
