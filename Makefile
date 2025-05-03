.DEFAULT_GOAL := help
.PHONY: devtest-helm

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
	kubectl apply --server-side --force-conflicts -f devtest-env/cnpg-1.25.1.yaml
	echo "Waiting for CNPG operator to be ready..."
	sleep 10 && kubectl wait --for=condition=Available deployment/cnpg-controller-manager -n cnpg-system --timeout=5m
	tilt up

test: ## Run tests using ginkgo. make dev should have been called ahead of this.
	@ORIG_URI=$$(kubectl -n vetchium-dev get secret postgres-app -o jsonpath='{.data.uri}' | base64 -d); \
	MOD_URI=$$(echo $$ORIG_URI | sed 's/postgres-rw.vetchium-dev/localhost/g'); \
	POSTGRES_URI=$$MOD_URI ginkgo -v ./dolores/...

seed: ## Seed the development database
	@ORIG_URI=$$(kubectl -n vetchium-dev get secret postgres-app -o jsonpath='{.data.uri}' | base64 -d); \
	MOD_URI=$$(echo $$ORIG_URI | sed 's/postgres-rw.vetchium-dev/localhost/g'); \
	cd dev-seed && POSTGRES_URI=$$MOD_URI go run .

lib: ## Build TypeSpec and install dependencies
	cd typespec && tsp compile . && npm run build && \
	cd ../harrypotter && npm install ../typespec && \
	cd ../ronweasly && npm install ../typespec

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
	docker buildx build -f sortinghat/Dockerfile \
		-t ghcr.io/vetchium/sortinghat:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--push .
	docker buildx build -f dev-seed/Dockerfile \
		-t ghcr.io/vetchium/dev-seed:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--push .

devtest-helm:
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
		--wait --timeout 10m

	@echo "=========================================================="
	@echo "Deployment complete! Use the following for distributed load testing:"
	@PGURI=$$(kubectl -n vetchium-devtest-$(VMUSER) get secret postgres-app -o jsonpath='{.data.uri}' 2>/dev/null | base64 -d 2>/dev/null | sed 's/postgres-rw.vetchium-devtest-'$(VMUSER)'/'$(VMADDR)'/g' 2>/dev/null || echo 'Error: Could not extract PostgreSQL URI')
	@echo "PostgreSQL URI: $$PGURI"
	@echo "VETCHIUM_API_SERVER_URL: http://$(VMADDR):8080"
	@echo "Mailpit URL: http://$(VMADDR):8025"
	@echo "
To run distributed load tests from a separate cluster:"
	@echo "VETCHIUM_API_SERVER_URL=http://$(VMADDR):8080 MAILPIT_URL=http://$(VMADDR):8025 PGURI=\"\$${PGURI}\" make k6-distributed"
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

# k6-distributed: Run distributed load tests with k6 using Kubernetes
# Parameter variables:
#   NUM_USERS/TOTAL_USERS - Total number of user accounts to create in the database (default: 1,000,000)
#   MAX_VUS - Maximum number of concurrent Virtual Users across all instances (default: 5,000)
#   INSTANCE_COUNT - Number of k6 instances to distribute the test across (default: 10)
#   TEST_DURATION - Duration of the test in seconds (default: 1800)
#   SETUP_PARALLELISM - Number of users to pre-authenticate during setup (default: 100)
#   VETCHIUM_API_SERVER_URL - (Required) Complete URL of the API server to test
#   MAILPIT_URL - (Required) Complete URL of the mailpit service
#   PGURI - (Required) PostgreSQL connection URI for the database
k6-distributed:
	@if [ -z "$(VETCHIUM_API_SERVER_URL)" ]; then \
		echo "Error: VETCHIUM_API_SERVER_URL environment variable is not set. This should be the complete URL of the API server."; \
		exit 1; \
	fi
	@if [ -z "$(MAILPIT_URL)" ]; then \
		echo "Error: MAILPIT_URL environment variable is not set. This should be the complete URL of the mailpit service."; \
		exit 1; \
	fi
	@if [ -z "$(PGURI)" ]; then \
		echo "Error: PGURI environment variable is not set. This should be the complete PostgreSQL connection URI."; \
		exit 1; \
	fi
	@echo "--- Checking connectivity to the backend services ---"
	@echo "Pinging API server at $(VETCHIUM_API_SERVER_URL)..."
	timeout 5 curl -s $(VETCHIUM_API_SERVER_URL)/healthz > /dev/null || echo "Warning: Could not reach API server. Make sure it's accessible."

	@echo "--- Creating k6 test namespace ---"
	# Generate a unique namespace for this test run
	$(eval K6_NAMESPACE := k6-loadtest-$(shell date +%Y%m%d-%H%M%S))
	kubectl create namespace $(K6_NAMESPACE)

	@echo "--- Creating database access configuration ---"
	# Store the provided PGURI in a secret for the test to use
	kubectl create secret generic postgres-credentials --from-literal=pguri="$(PGURI)" -n $(K6_NAMESPACE)

	@echo "--- Creating k6 distributed test resources ---"
	# Set variables for the test
	$(eval TEST_DURATION := $(or $(TEST_DURATION),1800))
	$(eval MAX_VUS := $(or $(MAX_VUS),5000))
	$(eval TOTAL_USERS := $(or $(NUM_USERS),1000000))
	$(eval INSTANCE_COUNT := $(or $(INSTANCE_COUNT),10))
	$(eval SETUP_PARALLELISM := $(or $(SETUP_PARALLELISM),100))

	# Apply variables to yaml template
	sed \
	    -e "s|\${VETCHIUM_API_SERVER_URL}|$(VETCHIUM_API_SERVER_URL)|g" \
	    -e "s|\${MAILPIT_URL}|$(MAILPIT_URL)|g" \
	    -e "s|\${TEST_DURATION}|$(TEST_DURATION)|g" \
	    -e "s|\${MAX_VUS}|$(MAX_VUS)|g" \
	    -e "s|\${TOTAL_USERS}|$(TOTAL_USERS)|g" \
	    -e "s|\${INSTANCE_COUNT}|$(INSTANCE_COUNT)|g" \
	    -e "s|\${SETUP_PARALLELISM}|$(SETUP_PARALLELISM)|g" \
	    ./neville/k6-distributed-updated.yaml | kubectl apply -f - -n $(K6_NAMESPACE)

	@echo "--- Copying test script to ConfigMap ---"
	kubectl create configmap k6-test-script --from-file=distributed_hub_scenario.js=neville/distributed_hub_scenario.js -n $(K6_NAMESPACE)

	@echo "--- Test started! ---"
	@echo "K6 test deployed in namespace: $(K6_NAMESPACE)"
	@echo "Monitor the test with: kubectl logs -f job/k6-distributed-test -n $(K6_NAMESPACE)"
	@echo "Individual worker logs: kubectl logs -f job/k6-worker -n $(K6_NAMESPACE) --selector=job-name=k6-worker"
	@echo "Target API server: $(VETCHIUM_API_SERVER_URL)"
	@echo "Target Mailpit server: $(MAILPIT_URL)"
	@echo "Clean up after test: kubectl delete namespace $(K6_NAMESPACE)"
