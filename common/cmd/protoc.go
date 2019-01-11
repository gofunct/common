// Copyright © 2019 Coleman Word <coleman.word@gofunct.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(protocCmd)
	protocCmd.PersistentFlags().StringSliceVarP(&plugins, "plugins", "p", []string{"gateway", "docs"}, "plugins to generate")
	protocCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output directory")

}

/*
--go_out=${GO_SOURCE_RELATIVE}plugins=grpc:$OUT_DIR"
		;;
		"gogo")
		GEN_STRING="--gogofast_out=${GO_SOURCE_RELATIVE}\
		Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,\
		Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
		Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,\
		Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
		Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,\
		Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,\
		plugins=grpc+embedded\

		-grpc_out=$OUT_DIR --${GEN_LANG}_out=$OUT_DIR --plugin=protoc-gen-grpc=`which grpc_${PLUGIN_LANG}_plugin
    mkdir -p $OUT_DIR/doc
        GEN_STRING="$GEN_STRING --doc_opt=$DOCS_FORMAT --doc_out=$OUT_DIR/doc"

        proto-include: -I /usr/local/include/ \
*/

// protocCmd represents the protoc command
var protocCmd = &cobra.Command{
	Use:   "protoc",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
