package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

func fatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func parseMF(pod string) (map[string]*dto.MetricFamily, error) {
	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(strings.NewReader(pod))
	if err != nil {
		return nil, err
	}
	return mf, nil
}

func readPods(reader *bufio.Reader, output chan string) error {
	defer close(output)
	first := true
	var buffer string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		if strings.HasPrefix(line, "# POD") {
			if first {
				first = false
			} else {
				output <- buffer
				buffer = ""
			}

		}
		buffer = buffer + line
	}
}

func main() {

	pods := make(chan string)

	reader := bufio.NewReader(os.Stdin)

	go readPods(reader, pods)

	for pod := range pods {
		mf, err := parseMF(pod)
		fatal(err)

		for k, v := range mf {
			fmt.Println("KEY: ", k)
			fmt.Println("VAL: ", v)
		}
	}
}
