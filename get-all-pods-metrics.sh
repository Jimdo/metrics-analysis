#!/bin/sh

# Set this to a prefix of a pod to skip reading in case of a previous error
skip_until_prefix=

set -euo pipefail # avoid surprises

context=$(kubectl config current-context)
metrics_output_path=metrics-$context

mkdir -p $metrics_output_path
skip=true

linkerd_pods=$(kubectl get pods --no-headers -llinkerd.io/control-plane-ns=linkerd -ojsonpath='{.items[*].metadata.name}')
for pod in $linkerd_pods
do
  case $pod in
    $skip_until_prefix*)
      skip=
      ;;
  esac

  if [ "$skip" == "true" ]; then
    continue
  fi

  echo "Reading po/$pod metrics"
  linkerd diagnostics proxy-metrics po/$pod >$metrics_output_path/$pod.metrics
done
