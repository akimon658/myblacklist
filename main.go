package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/idna"
)

func main() {
	templateFlag := flag.String("template", "", "Templates of filters. Put %s where a URL should be, use commas to specify multiple templates.")
	listFlag := flag.String("list", "", "List of URLs you want to block.")
	outputFlag := flag.String("output", "", "Location you want to save generated filter.")
	flag.Parse()

	if *templateFlag == "" {
		log.Fatal("Please specify 1 or more templates with the -template flag.")
	}
	templates := strings.Split(*templateFlag, ",")

	list, err := os.Open(*listFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer list.Close()

	filter, err := os.Create(*outputFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer filter.Close()

	w := bufio.NewWriter(filter)
	scanner := bufio.NewScanner(list)

	for scanner.Scan() {
		asciiUrl, err := idna.ToASCII(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		for i := range templates {
			w.WriteString(fmt.Sprintf(templates[i], asciiUrl))
			w.WriteString("\n")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	w.Flush()
}
