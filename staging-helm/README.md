# Vetchi Staging Environment

This Helm chart deploys the complete Vetchi application stack in a staging environment. It includes:

- PostgreSQL database cluster
- Backend services (Hermione and Granger)
- Frontend applications (Ron Weasly and Harry Potter)

## Prerequisites

- Linux or macOS system
- Sudo privileges
- Basic build tools (usually installed by default)

## Quick Start

Simply run:

```bash
make setup
```

This will:
1. Install Nix if not already installed
2. Set up a development environment with all required tools
3. Create a Kubernetes cluster using k3d
4. Build and load all Docker images
5. Deploy the application stack
6. Verify the installation

The setup process includes clear progress indicators and helpful messages at each step.

### First-Time Setup Notes

If this is your first time running the setup:

1. After Nix installation, you may need to:
   ```bash
   # Either restart your shell
   exec $SHELL

   # Or source the Nix profile
   source ~/.nix-profile/etc/profile.d/nix.sh
   ```

2. Then run setup again:
   ```bash
   make setup
   ```

### Installation Troubleshooting

If you encounter Nix installation issues:

1. **Missing nixbld group**: The setup will automatically create this for you
2. **Permission Issues**:
   ```bash
   # Ensure /nix directory has correct permissions
   sudo mkdir -p /nix
   sudo chown -R $USER:$USER /nix
   ```
3. **Shell Integration**:
   - The setup automatically adds Nix to both `.bashrc` and `.zshrc`
   - You may need to restart your shell after installation

4. **Other Issues**:
   - Ensure you have sudo privileges
   - Check system requirements at https://nixos.org/manual/nix/stable/installation/prerequisites.html
   - The setup will attempt both multi-user and single-user installations if needed

## Directory Structure

```
staging-helm/
├── Chart.yaml           # Helm chart metadata
├── values.yaml         # Default configuration values
├── Makefile           # Unified automation script
├── flake.nix          # Nix development environment
├── templates/         # Kubernetes manifests
│   ├── NOTES.txt     # Post-installation notes
│   ├── secrets.yaml  # Secret configurations
│   ├── hermione.yaml # Hermione backend service
│   ├── granger.yaml  # Granger backend service
│   ├── harrypotter.yaml # Harry Potter frontend
│   ├── ronweasly.yaml   # Ron Weasly frontend
│   └── postgres-cluster.yaml # Database cluster
├── Dockerfile.ronweasly    # Frontend Dockerfile
└── Dockerfile.harrypotter  # Frontend Dockerfile
```

## Configuration

### Global Settings

The following global settings can be configured in `values.yaml` or via `--set` flags:

```yaml
global:
  environment: staging
  domain: vetchi.org
```

### Component-specific Settings

Each component (Hermione, Granger, frontends) has its own configuration section in `values.yaml`. See the file for detailed options.

## Available Commands

- `make setup` - Complete setup (installs Nix, creates cluster, deploys applications)
- `make clean` - Remove everything (cluster, images, deployments)
- `make verify` - Check the status of all components
- `make uninstall` - Remove the application stack but keep the cluster

## Development Workflow

1. **Initial Setup**:
   ```bash
   make setup
   ```

2. **Make Changes**:
   - Edit Kubernetes manifests in `templates/`
   - Modify configuration in `values.yaml`
   - Update Docker images as needed

3. **Apply Changes**:
   ```bash
   # If you modified Docker images
   make build-images load-images

   # Update Helm release
   make install-chart
   ```

4. **Verify Changes**:
   ```bash
   make verify
   ```

## Accessing Services

After installation, services are available at:

- Harry Potter Frontend: http://localhost/harrypotter
- Ron Weasly Frontend: http://localhost/ronweasly
- Hermione API: http://localhost/api/hermione
- Granger API: http://localhost/api/granger

## Troubleshooting

1. **Check Pod Status**:
   ```bash
   kubectl get pods -n staging
   ```

2. **View Logs**:
   ```bash
   # Frontend logs
   kubectl logs -f -l app=harrypotter -n staging
   kubectl logs -f -l app=ronweasly -n staging

   # Backend logs
   kubectl logs -f -l app=hermione -n staging
   kubectl logs -f -l app=granger -n staging
   ```

3. **Common Issues**:
   - If this is your first time installing Nix, you may need to restart your shell
   - If pods are not starting, check resources with `kubectl describe pod <pod-name> -n staging`
   - For database issues, check PostgreSQL cluster status with `kubectl get postgresql -n staging`
   - For image pull issues, ensure images are properly loaded with `make load-images`

## Cleanup

To remove everything:

```bash
make clean
```

This will:
1. Uninstall the Helm chart
2. Delete the namespace
3. Remove the k3d cluster

## Contributing

1. Fork the repository
2. Create your feature branch
3. Make your changes
4. Submit a pull request

## License

This project is proprietary and confidential.
