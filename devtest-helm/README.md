To setup devtest on a fresh ubuntu VM, do:

```bash
su - <whatever-user-you-want>
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
$ # Setup context for kubectl
$ scp root@<public-ip-of-the-vm>:/etc/rancher/k3s/k3s.yaml k3s.yaml
$ export KUBECONFIG=$PWD/k3s.yaml
$ kubectl get pods -n vetchium-devtest-$USER

$ # Run k6 for load testing
vetchium $ VMUSER=<whatever-user-you-loggedin-the-vm> VMADDR=<public-ip-of-the-vm> make k6

$ # Port forward the services - Optionally if needed. This is not needed mostly, if the VMADDR services are reachable directly via VMADDR
vetchium $ VMUSER=<whatever-user-you-loggedin-the-vm> make port-forward-helm
```
