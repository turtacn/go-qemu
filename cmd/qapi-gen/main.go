package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/turtacn/go-qemu/qapi"
)

func main() {
	schemaDir := flag.String("schema-dir", "", "Path to QAPI schema directory")
	outDir := flag.String("out-dir", "", "Output directory for generated Go code")
	version := flag.String("version", "", "QEMU version string (e.g. 5.0.0)")
	pkg := flag.String("package", "", "Go package name for generated code")

	flag.Parse()

	if *schemaDir == "" || *outDir == "" || *version == "" || *pkg == "" {
		flag.Usage()
		os.Exit(1)
	}

	schema, err := qapi.ParseDir(*schemaDir, *version)
	if err != nil {
		log.Fatalf("Failed to parse schema in %s: %v", *schemaDir, err)
	}

	gen := qapi.NewGenerator(schema, *version, *pkg, *outDir)
	if err := gen.Run(); err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}

	fmt.Printf("Successfully generated Go code for QEMU %s in %s (package %s)\n", *version, *outDir, *pkg)
}
