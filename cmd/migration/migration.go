package main

import (
	"flag"
	"go-boilerplate/migrations"
)

func main() {
	flag.Parse()
	args := flag.Args()

	switch args[0] {
	case "migrate":
		migrations.Migrate(args[1:])
	}
}
