// Copyright 2018 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// file2byteslice is a dead simple tool to embed a file to Go.
package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	inputFilename  = flag.String("input", "", "input filename")
	outputFilename = flag.String("output", "", "output filename")
	packageName    = flag.String("package", "main", "package name")
	varName        = flag.String("var", "_", "variable name")
	compress       = flag.Bool("compress", false, "use gzip compression")
)

func write(w io.Writer, r io.Reader) error {
	if *compress {
		compressed := &bytes.Buffer{}
		cw, err := gzip.NewWriterLevel(compressed, gzip.BestCompression)
		if err != nil {
			return err
		}
		if _, err := io.Copy(cw, r); err != nil {
			return err
		}
		cw.Close()
		r = compressed
	}

	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, `// Code generated by file2byteslice. DO NOT EDIT.
// (gofmt is fine after generating)

package %s

var %s = []byte(%q)
`,
		*packageName, *varName, string(bs)); err != nil {
			return err
		}
	return nil
}

func run() error {
	var out io.Writer
	if *outputFilename != "" {
		f, err := os.Create(*outputFilename)
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	} else {
		out = os.Stdout
	}

	var in io.Reader
	if *inputFilename != "" {
		f, err := os.Open(*inputFilename)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	} else {
		in = os.Stdin
	}

	if err := write(out, in); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		panic(err)
	}
}