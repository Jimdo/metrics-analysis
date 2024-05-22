# Metrics Analysis

This tool can be used to analyze Prometheus metrics endpoints for high cardinality. If you want to contribute your own changes please take a look at the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## Usage

### Get a sample of all metrics data from the current cluster and namespace

This will store the metrics outputs in the directory `metrics-<current-context>`

```
./get-all-pods-metrics.sh
```

You can use the `skip_until_prefix` variable to skip reading already read
metrics in case of an error (e.g. a pod disappears while reading)

## Build the tool

```
make
```

### Usage

Use `./metrics-analysis -h` to see the help.

### List of Metrics (with summary)

```
./metrics-analysis -f metrics-.../podname.metrics -n
```

Output format: `<name> <label_count> (value_count, ...)`

### Show specific metric

```
./metrics-analysis -f metrics-.../podname.metrics -m response_latency_ms
```

### Summarize everything

```
./metrics-analysis -f metrics-.../podname.metrics
```
