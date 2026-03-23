package main

import (
	"github.com/ETEllis/teamcode/cmd"
	"github.com/ETEllis/teamcode/internal/logging"
)

func main() {
	defer logging.RecoverPanic("main", func() {
		logging.ErrorPersist("Application terminated due to unhandled panic")
	})

	cmd.Execute()
}
