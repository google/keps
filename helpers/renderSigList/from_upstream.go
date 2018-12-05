package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"os"

	"github.com/calebamiles/keps/helpers/renderSigList/internal/sigs"
)

func main() {
	sl, err := sigs.FetchUpstreamList()
	if err != nil {
		log.Fatalf("fetching upstream SIG information: %s", err)
	}

	renderedTemplate, err := ioutil.TempFile("", "upstream-sig-list")
	if err != nil {
		log.Fatalf("creating temp file for rendered template: %s", err)
	}

	err = sigs.RenderList(sl, renderedTemplate)
	if err != nil {
		os.Remove(renderedTemplate.Name())
		log.Fatalf("rendering SIG list template: %s", err)
	}

	err = renderedTemplate.Close()
	if err != nil {
		os.Remove(renderedTemplate.Name())
		log.Fatalf("closing rendered SIG list template: %s", err)
	}

	fmt.Printf("rendered template to: %s\n", renderedTemplate.Name())
	os.Exit(0)
}
