package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/Djarvur/parsewords"
)

const iterations = 100

func main() {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	text := string(bytes)

	var words []string

	started := time.Now()
	for ri := 0; ri < iterations; ri++ {
		words, err = parsewords.Shellwords(text)
		if err != nil {
			panic(err)
		}
	}
	spent := time.Since(started)

	fmt.Printf(
		"%d iterations done in %v: %d words found in %d bytes of input\n",
		iterations,
		spent,
		len(words),
		len(text),
	)
}
