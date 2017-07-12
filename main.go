package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/affinity226/ftpbeat/beater"
)

func main() {
	err := beat.Run("ftpbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
