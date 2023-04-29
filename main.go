package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/idna"
)

func main() {
	templateFlag := flag.String("template", "", "Path to template file used to generate a filter")
	listFlag := flag.String("list", "", "List of URLs you want to block")
	outputFlag := flag.String("output", "", "Location you want to save the generated filter")
	punycode := flag.Bool("punycode", false, "Enable Punycode encoding")
	flag.Parse()

	templateFile, err := os.Open(*templateFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer templateFile.Close()

	templScanner := bufio.NewScanner(templateFile)
	templates := make([]string, 0, 10)

	for templScanner.Scan() {
		templates = append(templates, templScanner.Text())
	}

	if err := templScanner.Err(); err != nil {
		log.Fatal(err)
	}

	list, err := os.Open(*listFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer list.Close()

	w := os.Stdout

	if *outputFlag != "" {
		filter, err := os.OpenFile(*outputFlag, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer filter.Close()

		w = filter
	}

	bw := bufio.NewWriter(w)
	scanner := bufio.NewScanner(list)

	for scanner.Scan() {
		s := scanner.Text()
		if *punycode {
			s, err = idna.ToASCII(s)
			if err != nil {
				log.Fatal(err)
			}
		}

		for i := range templates {
			bw.WriteString(fmt.Sprintf(templates[i], s))
			bw.WriteString("\n")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	bw.Flush()
}
