package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"

	"github.com/calebamiles/keps/pkg/sigs"
)

func main() {
	if len(os.Args) != 2 { // PROGNAME is the 0th argument
		log.Fatal("must only specifiy location to create SIG Directories as argument")
	}

	initLocation := os.Args[1]

	var errs *multierror.Error
	for _, sigName := range sigs.All() {
		errs = multierror.Append(errs, os.MkdirAll(filepath.Join(initLocation, sigName), os.ModePerm))
	}

	if errs.ErrorOrNil() != nil {
		log.Fatalf("creating SIG directories: %s", errs)
	}

	fmt.Printf("successfully created SIG directories at: %s\n", initLocation)
	os.Exit(0)
}
