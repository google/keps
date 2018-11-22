package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/calebamiles/keps/helpers/convert/internal/convert"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("must specifiy KEP location as argument")
	}

	kepLocation := os.Args[1]
	convertedLocation, err := convert.ToCurrent(kepLocation)
	if err != nil {
		log.Fatalf("converting KEP: %s", err)
	}

	fmt.Println("converted KEP content located at: " + convertedLocation)
	os.Exit(0)
}
