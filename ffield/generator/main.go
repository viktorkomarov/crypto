package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"
)

type Cfg struct {
	Degree      uint
	Path        string
	PackageName string
}

func parseConfig() Cfg {
	cfg := Cfg{}
	flag.UintVar(&cfg.Degree, "degree", 0, "gf(2^degree)")
	flag.StringVar(&cfg.Path, "path", "", "path to output")
	flag.StringVar(&cfg.PackageName, "package", "main", "name of generated file")
	flag.Parse()

	if cfg.Degree == 0 || cfg.Path == "" {
		fmt.Printf("usage of prorgamm: ffieldgen -path file -degree num -package name\n")
		os.Exit(1)
	}

	return cfg
}

//go:embed tables.tmpl
var templateTable string

type TemplateArgs struct {
	PackageName string
	Degree      uint
	SumOfTable  map[Pair]uint64
	MulOfTable  map[Pair]uint64
}

func main() {
	cfg := parseConfig()

	tableTmpl := template.New("tableTemplate")
	tableTmpl, err := tableTmpl.Parse(templateTable)
	if err != nil {
		log.Fatal(err)
	}

	gf, err := NewGF2(cfg.Degree)
	if err != nil {
		log.Fatalf("can't create gf %v\n", err)
	}

	err = tableTmpl.Execute(os.Stdout, TemplateArgs{
		PackageName: cfg.PackageName,
		Degree:      cfg.Degree,
		SumOfTable:  gf.generateSumTable(),
		MulOfTable:  gf.generateMulTable(),
	})
	if err != nil {
		log.Fatal(err)
	}
}
