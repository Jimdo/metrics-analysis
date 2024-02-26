package main

import (
	"bufio"
	"flag"
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
			output <- buffer
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

func readMetrics(pods chan string) map[string]map[string][]string {
	metrics := map[string]map[string][]string{}
	for pod := range pods {
		mf, err := parseMF(pod)
		fatal(err)

		for k, v := range mf {
			current := metrics[k]
			if current == nil {
				current = map[string][]string{}
			}

			labelValues := map[string][]string{}
			for _, m := range v.Metric {
				for _, l := range m.Label {
					labelValues[l.GetName()] = append(labelValues[l.GetName()], l.GetValue())
				}
			}

			for k, v := range labelValues {
				current[k] = append(current[k], v...)
			}

			metrics[k] = current
		}
	}

	return metrics
}

func unique(slice []string) []string {
	encountered := map[string]bool{}
	result := []string{}
	for v := range slice {
		if encountered[slice[v]] == true {
			continue
		}
		encountered[slice[v]] = true
		result = append(result, slice[v])
	}
	return result
}

func uniqMetrics(metrics map[string]map[string][]string) map[string]map[string][]string {
	uniq := map[string]map[string][]string{}
	for name, labels := range metrics {
		uniq[name] = map[string][]string{}
		for label_name, label_values := range labels {
			uniq[name][label_name] = unique(label_values)
		}
	}
	return uniq
}

func main() {
	metric_name := flag.String("m", "", "Metric name")
	list_metrics := flag.Bool("n", false, "List Metric names")
	flag.Parse()

	pods := make(chan string)

	reader := bufio.NewReader(os.Stdin)

	go readPods(reader, pods)

	metrics := readMetrics(pods)

	metrics = uniqMetrics(metrics)

	for name, labels := range metrics {
		if *metric_name != "" && *metric_name != name {
			continue
		}

		fmt.Printf("%s\n", name)
		if list_metrics != nil && *list_metrics {
			continue
		}
		for label_name, label_values := range labels {
			fmt.Printf("  %s:\n", label_name)
			for _, label_value := range label_values {
				fmt.Printf("    %s\n", label_value)
			}
		}
	}
}
