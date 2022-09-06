/*
* Copyright (c) 2015-2021 Alex Yatskov <alex@foosoft.net>
*
* Permission is hereby granted, free of charge, to any person obtaining a copy of
* this software and associated documentation files (the "Software"), to deal in
* the Software without restriction, including without limitation the rights to
* use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
* the Software, and to permit persons to whom the Software is furnished to do so,
* subject to the following conditions:
*
* The above copyright notice and this permission notice shall be included in all
* copies or substantial portions of the Software.
*
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
* IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
* FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
* COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
* IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
* CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"

	"github.com/russross/blackfriday"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] input output\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "https://foosoft.net/projects/md2vim/\n\n")
	fmt.Fprintf(os.Stderr, "Parameters:\n")
	flag.PrintDefaults()
}

func main() {
	cols := flag.Int("cols", defNumCols, "number of columns to use for layout")
	tabs := flag.Int("tabs", defTabSize, "tab width specified in number of spaces")
	notoc := flag.Bool("notoc", false, "do not generate table of contents for headings")
	norules := flag.Bool("norules", false, "do not generate horizontal rules above headings")
	pascal := flag.Bool("pascal", false, "use PascalCase for abbreviating tags")
	generate_tags := flag.Bool("generate-tags", false, "generate tags file")
	desc := flag.String("desc", "", "short description of the help file")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(2)
	}

	input, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Fatalf("unable to read from file: %s", args[0])
	}

	flags := 0
	if *notoc {
		flags |= flagNoToc
	}
	if *norules {
		flags |= flagNoRules
	}
	if *pascal {
		flags |= flagPascal
	}

	tags := make(map[string]struct{})
	renderer := VimDocRenderer(args[1], *desc, tags, *cols, *tabs, flags)
	extensions := blackfriday.EXTENSION_FENCED_CODE | blackfriday.EXTENSION_NO_INTRA_EMPHASIS | blackfriday.EXTENSION_SPACE_HEADERS | blackfriday.EXTENSION_HEADER_IDS
	output := blackfriday.Markdown(input, renderer, extensions)
	tagsSorted := make([]string, 0, len(tags))
    for k := range tags {
        tagsSorted = append(tagsSorted, k)
    }
	sort.Strings(tagsSorted)

	file, err := os.Create(args[1])
	if err != nil {
		log.Fatalf("unable to write to file: %s", args[1])
	}
	defer file.Close()

	if *generate_tags {

		dir := path.Dir(args[1])
		filename := path.Base(args[1])
		tags_file, err := os.Create(path.Join(dir, "tags"))

		if err != nil {
			log.Fatalf("unable to write to tag file")
		}
		defer tags_file.Close()

		for _, tag := range tagsSorted {
			tags_file.Write([]byte(fmt.Sprintf("%s\t%s\t/*%s*\n", tag, filename, tag)))
		}
	}


	if _, err := file.Write(output); err != nil {
		log.Fatal("unable to write output")
	}
}
