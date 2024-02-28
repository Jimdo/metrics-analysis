package main

import (
	"bytes"
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

func next(buf *[]byte, n int) []byte {
	ret := (*buf)[:n]
	(*buf) = (*buf)[n:]
	return ret
}

func readPods(file *string, output chan string) error {
	defer close(output)
	buf, _ := os.ReadFile(*file)
	for {
		needle := []byte("\n#\n# POD")
		split := bytes.Index(buf, needle)

		if split == -1 {
			break
		}

		split = split + 1 // we want to start with the '#'
		output <- string(next(&buf, split))
	}

	// no more needle then its the last
	output <- string(buf)
	return nil
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

func summariseOneMetric(labels map[string][]string) {
	fmt.Printf("%d (", len(labels))
	first := true
	for _, label_values := range labels {
		delimiter := ","
		if first {
			delimiter = ""
			first = false
		}
		fmt.Printf("%s%d", delimiter, len(label_values))
	}
	fmt.Printf(")\n")
}

func main() {
	metric_name := flag.String("m", "", "Metric name")
	list_metrics := flag.Bool("n", false, "List Metric names")
	list_labels := flag.Bool("l", false, "List Label names")
	select_label := flag.String("sl", "", "Filter for label name")
	file := flag.String("f", "", "Read from file")
	flag.Parse()

	pods := make(chan string)

	go readPods(file, pods)

	metrics := readMetrics(pods)
	metrics = uniqMetrics(metrics)

	for name, labels := range metrics {
		if *metric_name != "" && *metric_name != name {
			continue
		}

		fmt.Printf("%s ", name)
		summariseOneMetric(labels)

		if list_metrics != nil && *list_metrics {
			continue
		}
		for label_name, label_values := range labels {
			if select_label != nil && *select_label != "" {
				if label_name != *select_label {
					continue
				}
			}
			fmt.Printf("  %s: %d\n", label_name, len(label_values))
			if list_labels != nil && *list_labels {
				continue
			}
			for _, label_value := range label_values {
				fmt.Printf("    %s\n", label_value)
			}
		}
	}
}
