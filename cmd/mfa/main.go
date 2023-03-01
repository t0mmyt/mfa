package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/t0mmyt/mfa/storage"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	doCommand     = kingpin.Command("do", "do an mfa")
	doName        = doCommand.Arg("name", "name of mfa").Required().String()
	addCommand    = kingpin.Command("add", "add an mfa")
	addName       = addCommand.Arg("name", "name of mfa").Required().String()
	addUrl        = addCommand.Arg("url", "url of totp").Required().String()
	listCommand   = kingpin.Command("list", "list all stored keys")
	exportCommand = kingpin.Command("export", "export all keys as JSON")
	importCommand = kingpin.Command("import", "import all keys as JSON")
)

func main() {
	var keystore storage.Storage
	path, err := homedir.Expand("~/.local/share/mfa")
	if err != nil {
		log.Fatal(err)
	}

	keystore, err = storage.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := keystore.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	switch kingpin.Parse() {
	case addCommand.FullCommand():
		err := keystore.Set(*addName, *addUrl)
		if err != nil {
			log.Fatal(err)
		}
	case doCommand.FullCommand():
		key := keystore.Get(*doName)
		if key == nil {
			log.Fatalf("key %s not found", *doName)
		}
		code, err := key.Totp()
		if err != nil {
			log.Fatal(err)
		}
		println(code)
	case exportCommand.FullCommand():
		j, err := keystore.Export()
		if err != nil {
			log.Fatal(err)
		}
		_, err = fmt.Fprintf(os.Stdout, "%s", j)
		if err != nil {
			log.Fatal(err)
		}
	case importCommand.FullCommand():
		stdin, err := ioutil.ReadAll(os.Stdin)
		if len(stdin) == 0 {
			log.Fatal("Import reads from stdin")
		}
		if err != nil {
			log.Fatal(err)
		}
		err = keystore.Import(stdin)
		if err != nil {
			log.Fatal(err)
		}
	case listCommand.FullCommand():
		keys, err := keystore.List()
		if err != nil {
			log.Fatal(err)
		}
		for _, k := range keys {
			println(k)
		}
	}
}
