package yamldb

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func Load() *DB {
	f, err := os.Open("wtdb.yaml")
	if err != nil {
		log.Fatalln(fmt.Errorf("open db: %w", err))
	}
	defer f.Close()

	rawdb := &rawDB{}
	yaml.NewDecoder(f).Decode(rawdb)

	return rawdb.toDB()
}
