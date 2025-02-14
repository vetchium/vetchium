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
            curl

            # Core development tools
            go
            gopls
            gotools
            go-tools  # Additional Go tools like staticcheck
            docker

            # Node.js development
            nodejs_20    # LTS version
            nodePackages.npm
            nodePackages.pnpm
            nodePackages.yarn
            nodePackages.typescript
            nodePackages.typescript-language-server

            # Database tools
            postgresql_15

            # Additional tools
            jq      # JSON processing
            yq      # YAML processing
            git
          ];

          shellHook = ''
            # Ensure GOPATH is set
            export GOPATH="$HOME/go"
            export PATH="$GOPATH/bin:$PATH"

            # Create a directory for PID files
            mkdir -p .vetchi-pids

            # Function to check if a service is running
            is_running() {
              local pid_file=".vetchi-pids/$1.pid"
              if [ -f "$pid_file" ]; then
                local pid=$(cat "$pid_file")
                if kill -0 "$pid" 2>/dev/null; then
                  return 0
                fi
              fi
              return 1
            }

            # Function to start a Next.js app
            start_nextjs_app() {
              local app_name=$1
              local pid_file=".vetchi-pids/$app_name.pid"
              
              if ! is_running "$app_name"; then
                echo "📦 Starting $app_name..."
                cd "$app_name"
                if [ ! -d "node_modules" ]; then
                  echo "Installing dependencies for $app_name..."
                  npm install
                fi
                npm run dev > ../.vetchi-pids/$app_name.log 2>&1 &
                echo $! > "$pid_file"
                cd ..
                echo "✅ $app_name started (PID: $(cat $pid_file))"
              else
                echo "⚠️ $app_name is already running (PID: $(cat $pid_file))"
              fi
            }

            # Function to stop all services
            stop_services() {
              echo "Stopping all services..."
              for pid_file in .vetchi-pids/*.pid; do
                if [ -f "$pid_file" ]; then
                  local pid=$(cat "$pid_file")
                  local name=$(basename "$pid_file" .pid)
                  if kill -0 "$pid" 2>/dev/null; then
                    echo "Stopping $name (PID: $pid)..."
                    kill "$pid"
                  fi
                  rm "$pid_file"
                fi
              done
              
              # Stop tilt if running
              tilt down
            }

            # Function to show service status
            show_status() {
              echo "Service Status:"
              for pid_file in .vetchi-pids/*.pid; do
                if [ -f "$pid_file" ]; then
                  local pid=$(cat "$pid_file")
                  local name=$(basename "$pid_file" .pid)
                  if kill -0 "$pid" 2>/dev/null; then
                    echo "✅ $name is running (PID: $pid)"
                  else
                    echo "❌ $name is not running (stale PID file)"
                    rm "$pid_file"
                  fi
                fi
              done
            }

            # Register cleanup on shell exit
            trap stop_services EXIT

            # Start all services
            echo "🚀 Starting Vetchi development environment..."

            # Generate typespec libraries
            echo "Generating typespec libraries..."
            make lib

            # Start Kubernetes and Tilt
            echo "Starting Kubernetes environment..."
            make dev &
            echo $! > .vetchi-pids/tilt.pid

            # Start Next.js apps
            start_nextjs_app "harrypotter"
            start_nextjs_app "ronweasly"

            echo ""
            echo "💻 Development environment ready!"
            echo "Available commands:"
            echo "  status    - Show status of all services"
            echo "  stop      - Stop all services"
            echo ""
            echo "Logs are available in .vetchi-pids/*.log"
            echo ""
            show_status

            # Add convenience commands to the shell
            status() {
              show_status
            }

            stop() {
              stop_services
            }
          '';
        };
      }
    );
} 