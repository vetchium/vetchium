.PHONY: setup install-nix create-cluster build-images load-images install-chart uninstall clean verify check-deps

# Make setup the default target
.DEFAULT_GOAL := setup

SHELL := /bin/bash
.SHELLFLAGS := -ec

CLUSTER_NAME = vetchi-staging
NAMESPACE = staging
HELM_RELEASE = vetchi-staging
DOCKER_REGISTRY = psankar

# ANSI color codes
BLUE = \033[1;34m
GREEN = \033[1;32m
YELLOW = \033[1;33m
RED = \033[1;31m
RESET = \033[0m
CHECK = $(GREEN)✓$(RESET)
ARROW = $(BLUE)→$(RESET)

# Helper function to print section headers
define print_header
	printf "\n$(BLUE)════════════════════════════════════════════════════════════════$(RESET)\n"
	printf "$(BLUE)   %s$(RESET)\n" "$(1)"
	printf "$(BLUE)════════════════════════════════════════════════════════════════$(RESET)\n\n"
endef

# Helper function to print step progress
define print_step
	printf "$(ARROW) %s\n" "$(1)"
endef

# Helper function to print success
define print_success
	printf "$(CHECK) %s\n" "$(1)"
endef

# Check if command exists
check_cmd = $(shell command -v $(1) 2> /dev/null)

# Check if service is running
check_service = $(shell systemctl is-active $(1) 2>/dev/null)

# Check and start Docker
ensure-docker:
	$(call print_header,Checking Docker Status)
	$(call print_step,Verifying Docker installation...)
	command -v docker >/dev/null 2>&1 || { \
		printf "$(RED)Error: Docker is not installed. Installing...$(RESET)\n"; \
		curl -fsSL https://get.docker.com | sudo sh; \
	}
	$(call print_step,Checking Docker service...)
	if [ "$(call check_service,docker)" != "active" ]; then \
		printf "$(YELLOW)Docker service is not running. Starting...$(RESET)\n"; \
		sudo systemctl start docker || { \
			printf "$(RED)Failed to start Docker service. Please check:$(RESET)\n"; \
			printf "  • Run 'sudo systemctl status docker' for details\n"; \
			printf "  • Ensure your user is in the docker group: sudo usermod -aG docker $$USER\n"; \
			printf "  • You might need to log out and back in for group changes to take effect\n"; \
			exit 1; \
		}; \
		printf "Waiting for Docker to be ready...\n"; \
		sleep 5; \
	fi
	$(call print_success,Docker is running)
	$(call print_step,Checking user permissions...)
	if ! docker ps >/dev/null 2>&1; then \
		printf "$(YELLOW)Adding current user to docker group...$(RESET)\n"; \
		sudo usermod -aG docker $$USER; \
		printf "$(RED)Please log out and back in for group changes to take effect, then run make setup again$(RESET)\n"; \
		exit 1; \
	fi
	$(call print_success,Docker permissions verified)

# Help target
help:
	printf "Available targets:\n"
	printf "  setup         - Complete setup (installs Nix, creates cluster, deploys applications)\n"
	printf "  install-nix   - Install Nix package manager\n"
	printf "  create-cluster- Create k3d cluster\n"
	printf "  build-images  - Build all Docker images\n"
	printf "  load-images   - Load images into cluster\n"
	printf "  install-chart - Install Helm chart\n"
	printf "  verify        - Check installation status\n"
	printf "  uninstall     - Remove the application\n"
	printf "  clean         - Remove everything\n"

# Install Nix if not present
install-nix:
	$(call print_header,Checking and Installing Nix)
ifeq ($(call check_cmd,nix),)
	$(call print_step,Creating required system group 'nixbld'...)
	if ! getent group nixbld > /dev/null 2>&1; then \
		sudo groupadd -r nixbld || true; \
		for i in $$(seq 1 10); do \
			sudo useradd -c "Nix build user $$i" -d /var/empty -g nixbld -G nixbld -M -N -r -s $$(which nologin) nixbld$$i || true; \
		done; \
		$(call print_success,Created nixbld group and users); \
	else \
		$(call print_success,nixbld group already exists); \
	fi

	$(call print_step,Installing Nix...)
	if ! curl -L https://nixos.org/nix/install | sudo sh -s -- --daemon; then \
		printf "$(RED)Failed to install Nix. Trying single-user installation...$(RESET)\n"; \
		if ! curl -L https://nixos.org/nix/install | sh -s -- --no-daemon; then \
			printf "$(RED)Error: Both multi-user and single-user Nix installation failed.$(RESET)\n"; \
			printf "$(RED)Please check the following:$(RESET)\n"; \
			printf "  • Ensure you have sudo privileges\n"; \
			printf "  • Check if /nix directory exists and is writable\n"; \
			printf "  • Check system requirements at https://nixos.org/manual/nix/stable/installation/prerequisites.html\n"; \
			exit 1; \
		fi; \
	fi

	$(call print_step,Configuring Nix environment...)
	mkdir -p ~/.config/nix
	printf 'experimental-features = nix-command flakes\n' > ~/.config/nix/nix.conf
	if [ -f "$$HOME/.bashrc" ]; then \
		printf '\nif [ -e ~/.nix-profile/etc/profile.d/nix.sh ]; then\n  . ~/.nix-profile/etc/profile.d/nix.sh\nfi\n' >> ~/.bashrc; \
	fi
	if [ -f "$$HOME/.zshrc" ]; then \
		printf '\nif [ -e ~/.nix-profile/etc/profile.d/nix.sh ]; then\n  . ~/.nix-profile/etc/profile.d/nix.sh\nfi\n' >> ~/.zshrc; \
	fi

	$(call print_success,Nix installed successfully)
	printf "$(YELLOW)Important: You need to restart your shell or run:$(RESET)\n"
	printf "  source ~/.nix-profile/etc/profile.d/nix.sh\n"
else
	$(call print_success,Nix is already installed)
endif

# Check for required tools
check-deps:
	$(call print_header,Checking Dependencies)
	$(call print_step,Verifying Nix installation...)
	command -v nix >/dev/null 2>&1 || { printf "$(RED)Error: Nix is required but not installed. Run 'make install-nix' first.$(RESET)\n"; exit 1; }
	$(call print_success,All dependencies verified)

# Create k3d cluster
create-cluster:
	$(call print_header,Creating Kubernetes Cluster)
	$(call print_step,Creating k3d cluster '$(CLUSTER_NAME)'...)
	nix develop -c k3d cluster create $(CLUSTER_NAME) \
		--api-port 6550 \
		--port "80:80@loadbalancer" \
		--port "443:443@loadbalancer" \
		--agents 2
	$(call print_step,Adding PostgreSQL Helm repository...)
	nix develop -c helm repo add cnpg https://cloudnative-pg.github.io/charts
	nix develop -c helm repo update
	$(call print_success,Cluster created successfully)

# Build all Docker images
build-images:
	$(call print_header,Building Docker Images)
	$(call print_step,Building Hermione backend...)
	nix develop -c docker build -f ../api/Dockerfile-hermione -t $(DOCKER_REGISTRY)/hermione:latest ..
	$(call print_step,Building Granger backend...)
	nix develop -c docker build -f ../api/Dockerfile-granger -t $(DOCKER_REGISTRY)/granger:latest ..
	$(call print_step,Building Ron Weasly frontend...)
	nix develop -c docker build -f Dockerfile.ronweasly -t $(DOCKER_REGISTRY)/ronweasly:latest ..
	$(call print_step,Building Harry Potter frontend...)
	nix develop -c docker build -f Dockerfile.harrypotter -t $(DOCKER_REGISTRY)/harrypotter:latest ..
	$(call print_success,All images built successfully)

# Load images into k3d cluster
load-images: ensure-docker
	$(call print_header,Loading Images into Cluster)
	$(call print_step,Importing images into k3d...)
	nix develop -c k3d image import \
		$(DOCKER_REGISTRY)/hermione:latest \
		$(DOCKER_REGISTRY)/granger:latest \
		$(DOCKER_REGISTRY)/ronweasly:latest \
		$(DOCKER_REGISTRY)/harrypotter:latest \
		-c $(CLUSTER_NAME)
	$(call print_success,Images loaded successfully)

# Install Helm chart
install-chart:
	$(call print_header,Installing Helm Chart)
	$(call print_step,Creating namespace '$(NAMESPACE)'...)
	nix develop -c kubectl create namespace $(NAMESPACE) 2>/dev/null || true
	$(call print_step,Installing Helm chart...)
	nix develop -c helm install $(HELM_RELEASE) . \
		--namespace $(NAMESPACE) \
		--set global.domain=vetchi.local \
		--set global.environment=staging
	$(call print_success,Helm chart installed successfully)

# Verify installation
verify:
	$(call print_header,Verifying Installation)
	$(call print_step,Checking pod status...)
	nix develop -c kubectl get pods -n $(NAMESPACE)
	$(call print_step,Checking services...)
	nix develop -c kubectl get services -n $(NAMESPACE)
	$(call print_step,Checking PostgreSQL cluster...)
	nix develop -c kubectl get postgresql -n $(NAMESPACE)
	$(call print_success,Verification complete)

# Uninstall Helm chart
uninstall:
	$(call print_header,Uninstalling Application)
	$(call print_step,Removing Helm release...)
	nix develop -c helm uninstall $(HELM_RELEASE) -n $(NAMESPACE) 2>/dev/null || true
	$(call print_step,Removing namespace...)
	nix develop -c kubectl delete namespace $(NAMESPACE) 2>/dev/null || true
	$(call print_success,Uninstallation complete)

# Clean up everything
clean: uninstall
	$(call print_header,Cleaning Up Environment)
	$(call print_step,Deleting k3d cluster...)
	nix develop -c k3d cluster delete $(CLUSTER_NAME) 2>/dev/null || true
	$(call print_success,Cleanup complete)

# All-in-one setup
setup: install-nix check-deps create-cluster build-images load-images install-chart verify
	$(call print_header,Setup Complete)
	printf "$(GREEN)✨ Vetchi staging environment is ready!$(RESET)\n"
	printf "\n$(YELLOW)Access your applications at:$(RESET)\n"
	printf "  • Harry Potter Frontend: http://localhost/harrypotter\n"
	printf "  • Ron Weasly Frontend:  http://localhost/ronweasly\n"
	printf "  • Hermione API:         http://localhost/api/hermione\n"
	printf "  • Granger API:          http://localhost/api/granger\n"
	printf "\n$(YELLOW)To view logs and monitor the system:$(RESET)\n"
	printf "  • View all pods:        kubectl get pods -n $(NAMESPACE)\n"
	printf "  • View logs:            kubectl logs -f <pod-name> -n $(NAMESPACE)\n"
	printf "\n$(YELLOW)To clean up:$(RESET)\n"
	printf "  • Run:                  make clean\n"