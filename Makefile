.DEFAULT_GOAL := help

$(eval GIT_SHA=$(shell git rev-parse --short=18 HEAD))

help: ## Display this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z0-9_-]+:.*?## / { printf "  %-20s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

dev: ## Start development environment with Tilt and live reload
	tilt down
	time kubectl delete pvc -n vetchium-dev --all --ignore-not-found
	kubectl delete pv -n vetchium-dev --all --ignore-not-found
	kubectl delete namespace vetchium-dev --ignore-not-found --force --grace-period=0
	kubectl create namespace vetchium-dev
	helm repo add cnpg https://cloudnative-pg.github.io/charts
	helm upgrade --install cnpg \
		--namespace cnpg-system \
		--create-namespace \
		--version 0.22.1 \
		--wait --timeout 5m \
		cnpg/cloudnative-pg
	echo "Waiting for CNPG operator to be ready..."
	sleep 3 && kubectl wait --for=condition=Available deployment/cnpg-cloudnative-pg -n cnpg-system --timeout=5m
	tilt up

test: ## Run tests using ginkgo. make dev should have been called ahead of this.
	@ORIG_URI=$$(kubectl -n vetchium-dev get secret postgres-app -o jsonpath='{.data.uri}' | base64 -d); \
	MOD_URI=$$(echo $$ORIG_URI | sed 's/postgres-rw.vetchium-dev/localhost/g'); \
	POSTGRES_URI=$$MOD_URI ginkgo -v -r ./dolores/...

seed: ## Seed the development database
	@ORIG_URI=$$(kubectl -n vetchium-dev get secret postgres-app -o jsonpath='{.data.uri}' | base64 -d); \
	MOD_URI=$$(echo $$ORIG_URI | sed 's/postgres-rw.vetchium-dev/localhost/g'); \
	cd dev-seed && POSTGRES_URI=$$MOD_URI go run .

lib: ## Build TypeSpec and install dependencies
	cd typespec && tsp compile . && npm run build && \
	cd ../harrypotter && npm install ../typespec && \
	cd ../ronweasly && npm install ../typespec && \
	echo "TypeSpec Python models available at: typespec/sortinghat/"

docker: ## Build multi-platform Docker images
	docker buildx inspect multi-platform-builder >/dev/null 2>&1 || docker buildx create --name multi-platform-builder --platform=linux/amd64,linux/arm64 --use
	docker buildx build -f harrypotter/Dockerfile-optimized \
		-t ghcr.io/vetchium/harrypotter:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--build-arg API_ENDPOINT="http://hermione:8080" \
		.
	docker buildx build -f ronweasly/Dockerfile-optimized \
		-t ghcr.io/vetchium/ronweasly:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--build-arg API_ENDPOINT="http://hermione:8080" \
		.
	docker buildx build -f api/Dockerfile-hermione \
		-t ghcr.io/vetchium/hermione:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		.
	docker buildx build -f api/Dockerfile-granger \
		-t ghcr.io/vetchium/granger:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		.
	docker buildx build -f sqitch/Dockerfile \
		-t ghcr.io/vetchium/sqitch:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		sqitch
	docker buildx build -f sortinghat/Dockerfile.model-e5-base-v2 \
		-t ghcr.io/vetchium/sortinghat-model-e5-base-v2:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		.
	docker buildx build -f sortinghat/Dockerfile.model-bge-base-v1.5 \
		-t ghcr.io/vetchium/sortinghat-model-bge-base-v1.5:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		.
	docker buildx build -f sortinghat/Dockerfile.runtime \
		-t ghcr.io/vetchium/sortinghat:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		.
	docker buildx build -f dev-seed/Dockerfile \
		-t ghcr.io/vetchium/dev-seed:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		.

publish: ## Build multi-platform Docker images and publish them to the container registry
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "Error: There are uncommitted changes. Please commit them before publishing docker images."; \
		exit 1; \
	fi
	docker buildx inspect multi-platform-builder >/dev/null 2>&1 || docker buildx create --name multi-platform-builder --platform=linux/amd64,linux/arm64 --use
	docker buildx build -f harrypotter/Dockerfile-optimized \
		-t ghcr.io/vetchium/harrypotter:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--build-arg API_ENDPOINT="http://hermione:8080" \
		--push .
	docker buildx build -f ronweasly/Dockerfile-optimized \
		-t ghcr.io/vetchium/ronweasly:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--build-arg API_ENDPOINT="http://hermione:8080" \
		--push .
	docker buildx build -f api/Dockerfile-hermione \
		-t ghcr.io/vetchium/hermione:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--push .
	docker buildx build -f api/Dockerfile-granger \
		-t ghcr.io/vetchium/granger:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--push .
	docker buildx build -f sqitch/Dockerfile \
		-t ghcr.io/vetchium/sqitch:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--push sqitch
	docker buildx build -f sortinghat/Dockerfile.model-e5-base-v2 \
		-t ghcr.io/vetchium/sortinghat-model-e5-base-v2:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--push .
	docker buildx build -f sortinghat/Dockerfile.model-bge-base-v1.5 \
		-t ghcr.io/vetchium/sortinghat-model-bge-base-v1.5:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--push .
	docker buildx build -f sortinghat/Dockerfile.runtime \
		-t ghcr.io/vetchium/sortinghat:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--push .
	docker buildx build -f dev-seed/Dockerfile \
		-t ghcr.io/vetchium/dev-seed:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--push .

devtest: ## Deploy the development environment to a remote VM
	@if [ -z "$(VMUSER)" ]; then \
		echo "Error: VMUSER environment variable is not set."; \
		exit 1; \
	fi
	@if [ -z "$(VMADDR)" ]; then \
		echo "Error: VMADDR environment variable is not set. This should be the IP address where services will be accessible."; \
		exit 1; \
	fi
	helm uninstall vetchium-apps -n vetchium-devtest-$(VMUSER) || true # Optional: Uninstall previous app release first
	kubectl delete namespace vetchium-devtest-$(VMUSER) --ignore-not-found --force --grace-period=0
	kubectl create namespace vetchium-devtest-$(VMUSER)
	# Install/upgrade the environment chart first
	helm upgrade --install vetchium-env ./devtest-helm/vetchium-env-helm \
		--namespace vetchium-devtest-env \
		--create-namespace \
		--wait --timeout 10m
	# Wait for the CloudNativePG webhook service to be ready
	echo "Waiting for CloudNativePG webhook service to be ready..."
	kubectl wait --for=condition=Ready pod -l app.kubernetes.io/name=cloudnative-pg -n vetchium-devtest-env --timeout=5m
	# Now install/upgrade the applications chart
	helm upgrade --install vetchium-apps ./devtest-helm/vetchium-apps-helm \
		--namespace vetchium-devtest-$(VMUSER) \
		--create-namespace \
		--set global.vmaddr=$(VMADDR) \
		--set global.imageTag=$(GIT_SHA) \
		--wait --timeout 10m

	@echo "=========================================================="
	@echo "To run distributed load tests on a different cluster:"
	@PG_URI=$$(kubectl -n vetchium-devtest-$(VMUSER) get secret postgres-app -o jsonpath='{.data.uri}' | base64 --decode | sed 's/postgres-rw.vetchium-devtest-$(VMUSER)/$(VMADDR)/g'); \
	 echo "TOTAL_USERS=1000 TOTAL_PODS=5 VETCHIUM_API_SERVER_URL=http://$(VMADDR):8080 MAILPIT_URL=http://$(VMADDR):8025 PG_URI=$$PG_URI make k6"
	@echo "=========================================================="

port-forward-helm:
	# These are not needed mostly, if the VMADDR services are reachable directly via VMADDR
	# But added just for convenience
	@if [ -z "$(VMUSER)" ]; then \
		echo "Error: VMUSER environment variable is not set."; \
		exit 1; \
	fi
	pkill -9 -f "kubectl port-forward -n vetchium-devtest-$(VMUSER)" || true
	kubectl port-forward svc/harrypotter -n vetchium-devtest-$(VMUSER) 3001:80 &
	kubectl port-forward svc/ronweasly -n vetchium-devtest-$(VMUSER) 3002:80 &
	kubectl port-forward svc/mailpit-http -n vetchium-devtest-$(VMUSER) 8025:80 &
	kubectl port-forward svc/postgres-rw -n vetchium-devtest-$(VMUSER) 5432:5432 &
	kubectl port-forward svc/minio -n vetchium-devtest-$(VMUSER) 9000:9000 &
	kubectl port-forward svc/hermione -n vetchium-devtest-$(VMUSER) 8080:8080 &
	kubectl port-forward svc/grafana -n vetchium-devtest-env 3000:3000 &

# k6: Run distributed load tests with k6 using Kubernetes
# Parameter variables:
#   TOTAL_USERS - Total number of user accounts to create in the database (required)
#   INSTANCE_COUNT - Number of k6 instances to distribute the test across (required)
#   TEST_DURATION - Duration of the test in seconds (required)
#   SETUP_PARALLELISM - Number of users to pre-authenticate during setup (required)
#   VETCHIUM_API_SERVER_URL - Complete URL of the API server to test (required)
#   MAILPIT_URL - Complete URL of the mailpit service (required)
#   PG_URI - PostgreSQL connection URI for the database (required)
k6: ## Run distributed load tests with k6 using Kubernetes
	@if [ -z "$(VETCHIUM_API_SERVER_URL)" ]; then \
		echo "Error: VETCHIUM_API_SERVER_URL environment variable is not set. This should be the URL of the Vetchium API server."; \
		exit 1; \
	fi
	@if [ -z "$(MAILPIT_URL)" ]; then \
		echo "Error: MAILPIT_URL environment variable is not set. This should be the URL of the Mailpit server."; \
		exit 1; \
	fi
	@if [ -z "$(PG_URI)" ]; then \
		echo "Error: PG_URI environment variable is not set. This should be the complete PostgreSQL connection URI."; \
		exit 1; \
	fi
	# Set intelligent defaults and validate parameters
	$(eval TOTAL_USERS := $(or $(TOTAL_USERS),1000))
	$(eval TOTAL_PODS := $(or $(TOTAL_PODS),1))
	$(eval TEST_DURATION := $(or $(TEST_DURATION),660))  # 11 minutes (5 min peak + ramp up/down)

	# Validate parameters
	@if [ $(TOTAL_USERS) -lt 100 ]; then \
		echo "Error: TOTAL_USERS must be at least 100."; \
		exit 1; \
	fi

	@if [ $(TOTAL_PODS) -lt 1 ]; then \
		echo "Error: TOTAL_PODS must be at least 1."; \
		exit 1; \
	fi

	# Calculate users per pod and setup parallelism
	$(eval USERS_PER_POD := $(shell expr $(TOTAL_USERS) / $(TOTAL_PODS)))
	$(eval SETUP_PARALLELISM := $(or $(SETUP_PARALLELISM),$(shell expr $(USERS_PER_POD) / 10)))
	# Ensure minimum setup parallelism
	$(eval SETUP_PARALLELISM := $(shell if [ $(SETUP_PARALLELISM) -lt 10 ]; then echo 10; else echo $(SETUP_PARALLELISM); fi))

	@echo "--- Creating test users in the database ---"
	@echo "This may take some time for large user counts..."
	# Check if psql is installed
	which psql > /dev/null || { echo "Error: psql not found. Please install PostgreSQL client tools."; exit 1; }

	@echo "Executing SQL to create test users..."
	# Execute the SQL file with variable substitution using envsubst in a single network call
	@TOTAL_USERS="$(TOTAL_USERS)" envsubst '$$TOTAL_USERS' < ./neville/create_users.sql | \
	psql "$(PG_URI)" || { \
		echo "Error: Failed to create test users in the database. Aborting."; \
		exit 1; \
	}

	kubectl create namespace monitoring

	@echo "--- Helm Installing Prometheus+Grafana stack --- "
	helm upgrade --install kube-prometheus-stack \
		-f neville/kube-prometheus-stack-values.yaml \
		-n monitoring \
		prometheus-community/kube-prometheus-stack

	@echo "--- Creating k6 test namespace ---"
	# Generate a unique namespace for this test run
	$(eval K6_NAMESPACE := k6-loadtest-$(shell date +%Y%m%d-%H%M%S))
	kubectl create namespace $(K6_NAMESPACE)

	# Create the ConfigMap with the current script version
	kubectl create configmap k6-script --from-file=distributed_hub_scenario.js=neville/distributed_hub_scenario.js -n $(K6_NAMESPACE)

	@echo "--- Creating k6 distributed test resources ---"

	# Check if envsubst is installed
	which envsubst > /dev/null || { echo "Error: envsubst not found. Please install gettext package (brew install gettext on macOS or apt-get/yum install gettext on Linux)."; exit 1; }

	# Create k6 worker jobs for each pod
	@echo "Creating $(TOTAL_PODS) k6 worker pods..."
	@for i in $$(seq 0 $$(($(TOTAL_PODS) - 1))); do \
		echo "Creating k6 worker pod $$i of $(TOTAL_PODS)..."; \
		if [ $$i -gt 0 ]; then \
			DELAY=$$(( 5 + $$i * 3 + ($$RANDOM % 5) )); \
			echo "Waiting $$DELAY seconds before creating next pod..."; \
			sleep $$DELAY; \
		fi; \
		POD_INDEX=$$i \
		VETCHIUM_API_SERVER_URL=$(VETCHIUM_API_SERVER_URL) \
		MAILPIT_URL=$(MAILPIT_URL) \
		TEST_DURATION=$(TEST_DURATION) \
		TOTAL_USERS=$(TOTAL_USERS) \
		TOTAL_PODS=$(TOTAL_PODS) \
		USERS_PER_POD=$(USERS_PER_POD) \
		SETUP_PARALLELISM=$(SETUP_PARALLELISM) \
		PROMETHEUS_URL=$(PROMETHEUS_URL) \
		envsubst < ./neville/k6-job-template.yaml | kubectl apply -f - -n $(K6_NAMESPACE); \
	done

	# Verify that the pods are starting
	@echo "Waiting for k6 pods to start..."
	sleep 5
	kubectl get pods -n $(K6_NAMESPACE)

	@echo "--- Test started! ---"
	@echo "K6 test deployed in namespace: $(K6_NAMESPACE)"
	@echo ""
	@echo "To monitor test progress:"
	@echo "  kubectl get pods -n $(K6_NAMESPACE)"
	@echo ""
	@echo "To view logs from a specific pod:"
	@echo "  kubectl logs -f job/k6-worker-<POD_INDEX> -n $(K6_NAMESPACE)"
	@echo ""
	@echo "  Grafana credentials: admin / admin"
	@echo "  k6 dashboards are pre-installed in Grafana"
	@echo ""
	@echo "Target API server: $(VETCHIUM_API_SERVER_URL)"
	@echo "Target Mailpit server: $(MAILPIT_URL)"
	@echo ""
	@echo "Clean up after test: kubectl delete namespace $(K6_NAMESPACE)"
