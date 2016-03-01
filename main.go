package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/kusubooru/tags-diff/tags"
)

var (
	httpAddr   = flag.String("http", "localhost:8080", "HTTP listen address")
	staticPath = flag.String("static", "static/", "path to static files")
	certFile   = flag.String("tlscert", "", "TLS public key in PEM format.  Must be used together with -tlskey")
	keyFile    = flag.String("tlskey", "", "TLS private key in PEM format.  Must be used together with -tlscert")
	// Set after flag parsing based on certFile & keyFile.
	useTLS bool
)

const description = `Usage: tags-diff [options]
  A small utility that allows to find the differences between old tags and new
  tags which can be entered through a web interface.

Options:
`

func usage() {
	fmt.Fprintf(os.Stderr, description)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	flag.Usage = usage
	flag.Parse()
	useTLS = *certFile != "" && *keyFile != ""

	http.HandleFunc("/tags-diff", diff)

	fs := http.FileServer(http.Dir(*staticPath))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	if useTLS {
		if err := http.ListenAndServeTLS(*httpAddr, *certFile, *keyFile, nil); err != nil {
			log.Fatalf("Could not start listening (TLS) on %v: %v", *httpAddr, err)
		}
	} else {
		if err := http.ListenAndServe(*httpAddr, nil); err != nil {
			log.Fatalf("Could not start listening on %v: %v", *httpAddr, err)
		}
	}
}

func diff(w http.ResponseWriter, r *http.Request) {
	old := r.PostFormValue("old")
	new := r.PostFormValue("new")
	removed, added := tags.DiffFields(old, new)

	data := struct {
		Old     string
		New     string
		Removed []string
		Added   []string
	}{old, new, removed, added}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Print(err)
	}
}

// tmpl is the HTML template that drives the user interface.
var tmpl = template.Must(template.New("tmpl").Parse(`
<!DOCTYPE html><html><body>
	<style>
	ul,li {
	    list-style-type: none;
	}
	input, textarea {
		margin-bottom: 5px;
		display: block;
	}
	.removed {
		color: darkred;
	}
	.added {
		color: darkgreen;
	}
	</style>
	<form method="post" action="/tags-diff">
		<label for="old"><strong>Old Tags</strong></label>
		<textarea id="old" name="old" cols="60" rows="6">{{ .Old }}</textarea>
		<label for="new"><strong>New Tags</strong></label>
		<textarea id="new" name="new" cols="60" rows="6">{{ .New }}</textarea>
		<input type="submit" value="Compare">
	</form>
	<div id="diff">
	<samp>
	{{ range $r := .Removed }}
		<li><strong class="removed">---</strong> {{ $r }}</li>
	{{ end }}
	{{ range $a := .Added }}
		<li><strong class="added">+++</strong> {{ $a }}</li>
	{{ end }}
	</samp>
	</div>
</body></html>
`))
