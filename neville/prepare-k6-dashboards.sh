#!/bin/bash

set -e

echo "Preparing k6 dashboards for Grafana..."

# Create a temporary directory for dashboards
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

# Download k6 dashboards
echo "Downloading k6 dashboards..."

# Main k6 dashboard
curl -s -o "$TMP_DIR/k6-load-testing-results.json" https://raw.githubusercontent.com/grafana/k6/master/grafana/dashboards/k6-load-testing-results.json

# k6 Prometheus dashboard
curl -s -o "$TMP_DIR/k6-prometheus-dashboard.json" https://raw.githubusercontent.com/grafana/k6/master/grafana/dashboards/k6-prometheus.json

# Create ConfigMap for dashboards
echo "Creating ConfigMap for k6 dashboards..."
kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f -

# Create ConfigMap from downloaded dashboards
kubectl create configmap k6-dashboards \
  --from-file="k6-load-testing-results=$TMP_DIR/k6-load-testing-results.json" \
  --from-file="k6-prometheus=$TMP_DIR/k6-prometheus-dashboard.json" \
  -n monitoring --dry-run=client -o yaml | kubectl apply -f -

echo "k6 dashboards prepared successfully."
