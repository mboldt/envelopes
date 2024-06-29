package main

import (
	"encoding/csv"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Address struct {
	Name string
	Address string
	City string
	State string
	Zip string
}

func main() {
	r := csv.NewReader(os.Stdin)
	_, err := r.Read() // Ditch header
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	templateText, err := ioutil.ReadFile("envelope.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err := template.New("envelope").Parse(string(templateText))
	if err != nil {
		log.Fatal(err)
	}
	var html strings.Builder
	for _, record := range records {
		address := Address {
			record[1],
			record[2],
			record[3],
			record[4],
			record[5],
		}
		err = tmpl.Execute(&html, address)
		if err != nil {
			log.Fatal(err)
		}
	}
	wkhtmltopdf := exec.Command("wkhtmltopdf", "--page-width", "7.25in", "--page-height", "5.25in", "-", "e.pdf")
	stdin, err := wkhtmltopdf.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, html.String())
	}()
	err = wkhtmltopdf.Run()
	if err != nil {
		log.Fatal(err)
	}
}

