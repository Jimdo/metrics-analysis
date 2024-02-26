# Metrics Analysis

This tool can be used to analyze Prometheus metrics endpoints for high cardinality.

## Usage

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
