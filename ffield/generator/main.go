package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/viktorkomarov/crypto/bitset"
)

type Cfg struct {
	Degree      int
	Path        string
	PackageName string
}

func parseConfig() Cfg {
	cfg := Cfg{}
	flag.IntVar(&cfg.Degree, "degree", 0, "gf(2^degree)")
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
	Degree      int
	SumOfTable  map[string]*bitset.Set
}

func main() {
	cfg := parseConfig()

	tableTmpl := template.New("tableTemplate")
	tableTmpl, err := tableTmpl.Parse(templateTable)
	if err != nil {
		log.Fatal(err)
	}

	gf := NewGF2(cfg.Degree)
	sumTables := gf.GenerateSumTable(gf.Field())
	for _, val := range sumTables {
		fmt.Printf("%+v", val.Bits())
	}
	err = tableTmpl.Execute(os.Stdout, TemplateArgs{
		PackageName: cfg.PackageName,
		Degree:      cfg.Degree,
		SumOfTable:  sumTables,
	})
	if err != nil {
		log.Fatal(err)
	}
}
