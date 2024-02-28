#!/bin/sh

set -euo pipefail # avoid surprises

context=$(kubectl config current-context)
metrics_output_path=metrics-$context

mkdir -p $metrics_output_path

linkerd_pods=$(kubectl get pods --no-headers -llinkerd.io/control-plane-ns=linkerd -ojsonpath='{.items[*].metadata.name}')
for pod in $linkerd_pods
do
  echo "Reading po/$pod metrics"
  linkerd diagnostics proxy-metrics po/$pod >$metrics_output_path/$pod.metrics
done
