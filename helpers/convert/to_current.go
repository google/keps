package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync/atomic"

	log "github.com/sirupsen/logrus"

	"github.com/calebamiles/keps/helpers/convert/internal/convert"
)

func main() {
	if len(os.Args) != 2 { // PROGNAME is the 0th argument
		log.Fatal("must only specifiy KEP location as argument")
	}

	startLocation := os.Args[1]

	outputDir, err := ioutil.TempDir("", "kep-conversion-helper")
	if err != nil {
		log.Fatalf("creating output directory")
	}

	var filesToSkip = map[string]bool{
		"0004-cloud-provider-template.md": true,
	}

	var possibleKEPs uint32 = 0
	var convertedKEPs uint32 = 0

	err = filepath.Walk(startLocation, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if ext != ".md" {
			return nil
		}

		if filesToSkip[info.Name()] {
			log.Infof("skipping known file: %s", info.Name())
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("could not read content at: %s", path)
			return nil
		}

		if !bytes.Contains(content, []byte("---")) {
			// document unlikely to contain KEP metadata
			return nil
		}

		atomic.AddUint32(&possibleKEPs, 1)

		convertedKEPLocation, err := convert.ToCurrent(outputDir, path)
		if err != nil {
			log.Errorf("failed to convert possible KEP at: %s with error: %s", path, err)
			return nil // keep going and try to convert the next KEP
		}

		atomic.AddUint32(&convertedKEPs, 1)
		log.Infof("converted KEP saved at: %s", convertedKEPLocation)

		return nil
	})

	fmt.Println("converted KEP content located at: " + outputDir)
	fmt.Printf("converted %d out of %d KEPs: %f percent success rate\n", atomic.LoadUint32(&convertedKEPs), atomic.LoadUint32(&possibleKEPs), 100.0*float32(atomic.LoadUint32(&convertedKEPs))/float32(atomic.LoadUint32(&possibleKEPs)))
	os.Exit(0)
}
