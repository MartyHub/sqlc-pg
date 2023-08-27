package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/MartyHub/sqlc-pg/internal"
	"github.com/MartyHub/sqlc-pg/plugin"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	req, err := parseInput()
	if err != nil {
		return err
	}

	generator, err := internal.New(req)
	if err != nil {
		return err
	}

	rep, err := generator.Generate()
	if err != nil {
		return err
	}

	return writeOutput(rep)
}

func parseInput() (*plugin.CodeGenRequest, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	result := new(plugin.CodeGenRequest)

	if err = result.UnmarshalVT(data); err != nil {
		return nil, err
	}

	return result, err
}

func writeOutput(rep *plugin.CodeGenResponse) error {
	data, err := rep.MarshalVT()
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(os.Stdout)
	if _, err = writer.Write(data); err != nil {
		return err
	}

	return writer.Flush()
}
