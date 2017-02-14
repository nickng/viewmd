// Command viewmd renders a Markdown file and displays the output in browser.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/browser"
	"github.com/russross/blackfriday"
)

var (
	help bool
	keep bool
	wait int
)

func init() {
	flag.BoolVar(&help, "help", false, "Show this help message")
	flag.BoolVar(&keep, "keep", false, "Keep the generated HTML file")
	flag.IntVar(&wait, "wait", 3, "Number of seconds to wait "+
		"before deleting generated files (ignored if -keep is set)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"%s renders a Markdown file and displays the output in browser.\n\n"+
				"Usage:\n\t%s [options] input.md...\nOptions:\n",
			os.Args[0], os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.Parse()
}

func main() {
	if flag.NArg() == 0 || help {
		flag.Usage()
	}

	tempdir, err := ioutil.TempDir("", "viewmd")
	if err != nil {
		log.Fatal("Cannot create temporary output directory:", err)
	}

	for i, arg := range flag.Args() {
		b, err := ioutil.ReadFile(arg)
		if err != nil {
			log.Fatal("Cannot read input markdown file:", err)
		}
		unsafe := blackfriday.MarkdownCommon(b)
		html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
		outfile := path.Join(tempdir, fmt.Sprintf("output%d.html", i))
		if err := ioutil.WriteFile(outfile, html, 0600); err != nil {
			log.Fatal("Cannot write output to file: ", err)
		}
		log.Println("Output written to", outfile)
		if err := browser.OpenFile(outfile); err != nil {
			log.Fatal("Cannot open file: ", err)
		}
	}

	if !keep {
		// Allow time for file to be read before removing
		time.Sleep(time.Duration(wait) * time.Second)
		os.RemoveAll(tempdir)
	}
}
