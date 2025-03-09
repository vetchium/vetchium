.PHONY: setup install-nix create-cluster build-images load-images install-chart uninstall clean verify check-deps

CLUSTER_NAME = vetchi-staging
NAMESPACE = staging
HELM_RELEASE = vetchi-staging
DOCKER_REGISTRY = psankar

# ANSI color codes
BLUE := \033[1;34m
GREEN := \033[1;32m
YELLOW := \033[1;33m
RED := \033[1;31m
RESET := \033[0m
CHECK := $(GREEN)✓$(RESET)
ARROW := $(BLUE)→$(RESET)

# Helper function to print section headers
define print_header
	@echo "\n$(BLUE)════════════════════════════════════════════════════════════════$(RESET)"
	@echo "$(BLUE)   $1$(RESET)"
	@echo "$(BLUE)════════════════════════════════════════════════════════════════$(RESET)\n"
endef

# Helper function to print step progress
define print_step
	@echo "$(ARROW) $1"
endef

# Helper function to print success
define print_success
	@echo "$(CHECK) $1"
endef

# Check if command exists
check_cmd = $(shell command -v $(1) 2> /dev/null)

# Install Nix if not present
install-nix:
	$(call print_header,Checking and Installing Nix)
ifeq ($(call check_cmd,nix),)
	$(call print_step,Creating required system group 'nixbld'...)
	@if ! getent group nixbld > /dev/null 2>&1; then \
		sudo groupadd -r nixbld; \
		for i in $$(seq 1 10); do \
			sudo useradd -c "Nix build user $$i" -d /var/empty -g nixbld -G nixbld -M -N -r -s "$(shell which nologin)" nixbld$$i; \
		done; \
		$(call print_success,Created nixbld group and users); \
	else \
		$(call print_success,nixbld group already exists); \
	fi

	$(call print_step,Installing Nix...)
	@if ! sh <(curl -L https://nixos.org/nix/install) --daemon; then \
		echo "$(RED)Failed to install Nix. Trying single-user installation...$(RESET)"; \
		if ! sh <(curl -L https://nixos.org/nix/install) --no-daemon; then \
			echo "$(RED)Error: Both multi-user and single-user Nix installation failed.$(RESET)"; \
			echo "$(RED)Please check the following:$(RESET)"; \
			echo "  • Ensure you have sudo privileges"; \
			echo "  • Check if /nix directory exists and is writable"; \
			echo "  • Check system requirements at https://nixos.org/manual/nix/stable/installation/prerequisites.html"; \
			exit 1; \
		fi; \
	fi

	$(call print_step,Configuring Nix environment...)
	@{ \
		echo 'if [ -e ~/.nix-profile/etc/profile.d/nix.sh ]; then'; \
		echo '  . ~/.nix-profile/etc/profile.d/nix.sh'; \
		echo 'fi'; \
	} >> ~/.bashrc
	@{ \
		echo 'if [ -e ~/.nix-profile/etc/profile.d/nix.sh ]; then'; \
		echo '  . ~/.nix-profile/etc/profile.d/nix.sh'; \
		echo 'fi'; \
	} >> ~/.zshrc

	$(call print_success,Nix installed successfully)
	@echo "$(YELLOW)Important: You need to restart your shell or run:$(RESET)"
	@echo "  source ~/.nix-profile/etc/profile.d/nix.sh"
else
	$(call print_success,Nix is already installed)
endif

# Check for required tools
check-deps:
	$(call print_header,Checking Dependencies)
	$(call print_step,Verifying Nix installation...)
	@command -v nix >/dev/null 2>&1 || { echo "$(RED)Error: Nix is required but not installed. Run 'make install-nix' first.$(RESET)"; exit 1; }
	$(call print_success,All dependencies verified)

# Create k3d cluster
create-cluster:
	$(call print_header,Creating Kubernetes Cluster)
	$(call print_step,Creating k3d cluster '$(CLUSTER_NAME)'...)
	@nix develop -c k3d cluster create $(CLUSTER_NAME) \
		--api-port 6550 \
		--port "80:80@loadbalancer" \
		--port "443:443@loadbalancer" \
		--agents 2
	$(call print_step,Adding PostgreSQL Helm repository...)
	@nix develop -c helm repo add cnpg https://cloudnative-pg.github.io/charts
	@nix develop -c helm repo update
	$(call print_success,Cluster created successfully)

# Build all Docker images
build-images:
	$(call print_header,Building Docker Images)
	$(call print_step,Building Hermione backend...)
	@nix develop -c docker build -f ../api/Dockerfile-hermione -t $(DOCKER_REGISTRY)/hermione:latest ..
	$(call print_step,Building Granger backend...)
	@nix develop -c docker build -f ../api/Dockerfile-granger -t $(DOCKER_REGISTRY)/granger:latest ..
	$(call print_step,Building Ron Weasly frontend...)
	@nix develop -c docker build -f Dockerfile.ronweasly -t $(DOCKER_REGISTRY)/ronweasly:latest ..
	$(call print_step,Building Harry Potter frontend...)
	@nix develop -c docker build -f Dockerfile.harrypotter -t $(DOCKER_REGISTRY)/harrypotter:latest ..
	$(call print_success,All images built successfully)

# Load images into k3d cluster
load-images:
	$(call print_header,Loading Images into Cluster)
	$(call print_step,Importing images into k3d...)
	@nix develop -c k3d image import \
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
	@nix develop -c kubectl create namespace $(NAMESPACE) 2>/dev/null || true
	$(call print_step,Installing Helm chart...)
	@nix develop -c helm install $(HELM_RELEASE) . \
		--namespace $(NAMESPACE) \
		--set global.domain=vetchi.local \
		--set global.environment=staging
	$(call print_success,Helm chart installed successfully)

# Verify installation
verify:
	$(call print_header,Verifying Installation)
	$(call print_step,Checking pod status...)
	@nix develop -c kubectl get pods -n $(NAMESPACE)
	$(call print_step,Checking services...)
	@nix develop -c kubectl get services -n $(NAMESPACE)
	$(call print_step,Checking PostgreSQL cluster...)
	@nix develop -c kubectl get postgresql -n $(NAMESPACE)
	$(call print_success,Verification complete)

# Uninstall Helm chart
uninstall:
	$(call print_header,Uninstalling Application)
	$(call print_step,Removing Helm release...)
	@nix develop -c helm uninstall $(HELM_RELEASE) -n $(NAMESPACE) 2>/dev/null || true
	$(call print_step,Removing namespace...)
	@nix develop -c kubectl delete namespace $(NAMESPACE) 2>/dev/null || true
	$(call print_success,Uninstallation complete)

# Clean up everything
clean: uninstall
	$(call print_header,Cleaning Up Environment)
	$(call print_step,Deleting k3d cluster...)
	@nix develop -c k3d cluster delete $(CLUSTER_NAME) 2>/dev/null || true
	$(call print_success,Cleanup complete)

# All-in-one setup
setup: install-nix check-deps create-cluster build-images load-images install-chart verify
	$(call print_header,Setup Complete)
	@echo "$(GREEN)✨ Vetchi staging environment is ready!$(RESET)"
	@echo "\n$(YELLOW)Access your applications at:$(RESET)"
	@echo "  • Harry Potter Frontend: http://localhost/harrypotter"
	@echo "  • Ron Weasly Frontend:  http://localhost/ronweasly"
	@echo "  • Hermione API:         http://localhost/api/hermione"
	@echo "  • Granger API:          http://localhost/api/granger"
	@echo "\n$(YELLOW)To view logs and monitor the system:$(RESET)"
	@echo "  • View all pods:        kubectl get pods -n $(NAMESPACE)"
	@echo "  • View logs:            kubectl logs -f <pod-name> -n $(NAMESPACE)"
	@echo "\n$(YELLOW)To clean up:$(RESET)"
	@echo "  • Run:                  make clean"