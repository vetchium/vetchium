{
  description = "Vetchi Staging Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # Create a script to ensure docker daemon is running
        dockerEnsureScript = pkgs.writeShellScriptBin "ensure-docker" ''
          if ! docker info >/dev/null 2>&1; then
            echo "Starting Docker daemon..."
            sudo ${pkgs.docker}/bin/dockerd &
            # Wait for Docker to be ready
            for i in {1..30}; do
              if docker info >/dev/null 2>&1; then
                break
              fi
              echo "Waiting for Docker to be ready... $i/30"
              sleep 1
            done
          fi
        '';
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Core tools
            kubectl
            kubernetes-helm
            k3d

            # Docker and related tools
            docker
            docker-compose
            dockerEnsureScript

            # Build tools
            gnumake
            curl
            jq
          ];

          # Set up environment
          shellHook = ''
            # Ensure Docker socket exists and is accessible
            if [ ! -S /var/run/docker.sock ]; then
              echo "$(tput setaf 3)Docker socket not found. Setting up Docker...$(tput sgr0)"
              sudo mkdir -p /var/run
              sudo ln -sf $DOCKER_HOST /var/run/docker.sock
              ensure-docker
            fi

            # Print available commands
            echo "$(tput setaf 2)Vetchi Staging Development Environment$(tput sgr0)"
            echo "$(tput setaf 6)Available commands:$(tput sgr0)"
            echo "  make setup      - Complete setup"
            echo "  make create-cluster - Create k3d cluster"
            echo "  make build-images   - Build Docker images"
            echo "  make load-images    - Load images into cluster"
            echo "  make install-chart  - Install Helm chart"
            echo "  make verify     - Verify installation"
            echo "  make uninstall  - Uninstall chart"
            echo "  make clean      - Clean everything"
          '';

          # Set environment variables
          DOCKER_BUILDKIT = "1";
          COMPOSE_DOCKER_CLI_BUILD = "1";
        };
      }
    );
} 