# Installation Guide

## Prerequisites (Manual Installation)

1. Install Docker:

```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
# Log out and log back in for group changes to take effect
```

2. Install kubectl:

```bash
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
```

3. Install Helm and add the CloudnativePG repo

```bash
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
helm repo add cnpg https://cloudnative-pg.github.io/charts
helm repo update
```

4. Install k3d (lightweight Kubernetes):

```bash
curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
```

## Prerequisites (Via cloud-init)

The above packages can also be installed via cloud-init in supported VM providers, using the [cloud-init.yaml](cloud-init.yaml) file

## Setup Kubernetes Cluster

Create a k3d cluster (or) set the kubectx for any k8s cluster:

```bash
k3d cluster create vetchi-staging \
  --api-port 6550 \
  --port "80:80@loadbalancer" \
  --port "443:443@loadbalancer" \
  --agents 2
```

## Build and Load Docker Images

1. Build the Docker images:

```bash
# Build backend images
docker build -f api/Dockerfile-hermione -t psankar/hermione:latest .
docker build -f api/Dockerfile-granger -t psankar/granger:latest .
docker build -f sqitch/Dockerfile -t psankar/sqitch:latest sqitch

# Build frontend images
docker build -f staging-helm/Dockerfile-ronweasly -t psankar/ronweasly:latest .
docker build -f staging-helm/Dockerfile-harrypotter -t psankar/harrypotter:latest .

# Import images into k3d cluster
k3d image import psankar/hermione:latest psankar/granger:latest \
  psankar/ronweasly:latest psankar/harrypotter:latest \
  psankar/sqitch:latest -c vetchi-staging
```

## Install the Helm Chart

1. Create namespace:

```bash
kubectl create namespace staging
```

2. Install the chart:

```bash
helm install vetchi-staging ./staging-helm \
  --namespace staging \
  --set global.domain=vetchi.local \
  --set global.environment=staging
```

## Verify Installation

1. Check if all pods are running:

```bash
kubectl get pods -n staging
```

2. Check services:

```bash
kubectl get services -n staging
```

3. Access the applications:

- Harry Potter Frontend: http://localhost/harrypotter
- Ron Weasly Frontend: http://localhost/ronweasly
- Hermione API: http://localhost/api/hermione
- Granger API: http://localhost/api/granger

## Troubleshooting

1. Check pod logs:

```bash
kubectl logs -f -l app=harrypotter -n staging
kubectl logs -f -l app=ronweasly -n staging
kubectl logs -f -l app=hermione -n staging
kubectl logs -f -l app=granger -n staging
```

2. Check pod status:

```bash
kubectl describe pods -n staging
```

3. Check PostgreSQL cluster status:

```bash
kubectl get postgresql -n staging
kubectl get pods -l postgresql -n staging
```

## Uninstall

To remove the installation:

```bash
helm uninstall vetchi-staging -n staging
kubectl delete namespace staging
```

To delete the k3d cluster:

```bash
k3d cluster delete vetchi-staging
```
