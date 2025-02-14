{
  description = "Vetchi Development Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Kubernetes development tools
            kubectl
            tilt
            kubernetes-helm
            k9s  # Terminal UI for Kubernetes

            # Core development tools
            go_1_21
            gopls
            gotools
            go-tools  # Additional Go tools like staticcheck
            docker

            # Database tools
            postgresql_15
            sqitch  # For database migrations

            # Additional tools
            jq      # JSON processing
            yq      # YAML processing
            git
          ];

          shellHook = ''
            # Ensure GOPATH is set
            export GOPATH="$HOME/go"
            export PATH="$GOPATH/bin:$PATH"

            echo "🚀 Vetchi development environment loaded!"
            echo "Run 'tilt up' to start your local development environment"
          '';
        };
      }
    );
} 