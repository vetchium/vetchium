To setup devtest on a fresh ubuntu VM, do:

```bash
snap install docker
snap install helm --classic
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--tls-san <public-ip-of-the-vm>" sh -
git clone https://github.com/vetchium/vetchium.git
cd ~/vetchium/devtest-helm/vetchium-apps-helm
helm dependency update .
cd ~/vetchium/devtest-helm/vetchium-env-helm
helm dependency update .
cd ~/vetchium
make devtest-helm
```

To get the kubectl access and port forward the services, on your developer laptop for the above service, do:
```bash
scp root@<public-ip-of-the-vm>:/etc/rancher/k3s/k3s.yaml k3s.yaml
export KUBECONFIG=$PWD/k3s.yaml
vetchium $ kubectl get pods -n vetchium-devtest-$USER
vetchium $ make port-forward-helm
vetchium $ make k6
```
