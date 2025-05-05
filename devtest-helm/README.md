To setup helm based backend pods on a fresh ubuntu VM, do:
```bash
apt update && apt install make
snap install docker
snap install helm --classic

# Ensure that the firewall is disabled and stopped
systemctl disable ufw
systemctl stop ufw
ufw status

# Create a non-root user and switch to it
useradd -m -g users -G sudo -s /bin/bash <whatever-user-you-want>
passwd <whatever-user-you-want> # Set the password
su - <whatever-user-you-want>

# Install k3s and setup the environment
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--tls-san <public-ip-of-the-vm> --write-kubeconfig-mode 644" sh -
git clone https://github.com/vetchium/vetchium.git
cd ~/vetchium/devtest-helm/vetchium-apps-helm
helm dependency update .
cd ~/vetchium/devtest-helm/vetchium-env-helm
helm dependency update .
cd ~/vetchium
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml; # This is needed for helm. kubectl will work even otherwise via k3s init script
VMUSER=<whatever-user-you-want> VMADDR=<public-ip-of-the-vm> make devtest
kubectl get pods -n vetchium-devtest-<whatever-user-you-want>

# Make sure postgres-rw is LoadBalancer
kubectl patch service postgres-rw \
  -p '{"spec": {"type": "LoadBalancer"}}' \
  -n vetchium-devtest-<whatever-user-you-want>

kubectl get svc -n vetchium-devtest-<whatever-user-you-want>
```

To port forward the services - Optionally if needed. This is not needed mostly.
```bash
vetchium $ VMUSER=<whatever-user-you-loggedin-the-vm> make port-forward-helm
```

To run k6 for load testing, do, from your laptop, with a different k8s context:
```bash
vetchium $ TOTAL_USERS=100 \
  TOTAL_PODS=5 \
  VETCHIUM_API_SERVER_URL=http://<public-ip-of-the-vm>:8080 \
  MAILPIT_URL=http://<public-ip-of-the-vm>:8025 \
  PG_URI=postgresql://app:<password>@<public-ip-of-the-vm>:5432/app \
  make k6
```

Note: When running tests from a separate cluster, make sure network connectivity and firewall rules allow access from your test cluster to the backend services (API server, Mailpit, and PostgreSQL database)
