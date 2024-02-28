# Metrics Analysis

This tool can be used to analyze Prometheus metrics endpoints for high cardinality.

## Usage

### Get a sample of all metrics data from the current cluster and namespace

This will store the metrics outputs in the directory `metrics-<current-context>`

```
./get-all-pods-metrics.sh
```

### List of Metrics (with summary)

```
linkerd diagnostics proxy-metrics deploy/<service> | go run metrics.go -n
```

Output format: `<name> <label_count> (value_count, ...)`

### Show specific metric

```
linkerd diagnostics proxy-metrics deploy/<service> | go run metrics.go -m response_latency_ms
```

### Summarize everything

```
linkerd diagnostics proxy-metrics deploy/<service> | go run metrics.go
```
