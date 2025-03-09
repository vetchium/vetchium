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
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            docker
            kubectl
            kubernetes-helm
            k3d
            gnumake
          ];

          shellHook = ''
            echo "Vetchi Staging Development Environment"
            echo "Available commands:"
            echo "  make setup      - Complete setup"
            echo "  make install-deps   - Install dependencies"
            echo "  make create-cluster - Create k3d cluster"
            echo "  make build-images   - Build Docker images"
            echo "  make load-images    - Load images into cluster"
            echo "  make install-chart  - Install Helm chart"
            echo "  make verify     - Verify installation"
            echo "  make uninstall  - Uninstall chart"
            echo "  make clean      - Clean everything"
          '';
        };
      }
    );
} 