.DEFAULT_GOAL := help

$(eval GIT_SHA=$(shell git rev-parse --short=18 HEAD))

help: ## Display this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z0-9_-]+:.*?## / { printf "  %-20s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

dev: ## Start development environment with Tilt and live reload
	tilt down
	time kubectl delete pvc -n vetchidev --all --ignore-not-found
	kubectl delete pv -n vetchidev --all --ignore-not-found
	kubectl delete namespace vetchidev --ignore-not-found --force --grace-period=0
	kubectl create namespace vetchidev
	kubectl apply --server-side --force-conflicts -f devtest-env/cnpg-1.25.1.yaml
	echo "Waiting for CNPG operator to be ready..."
	sleep 10 && kubectl wait --for=condition=Available deployment/cnpg-controller-manager -n cnpg-system --timeout=5m
	tilt up

test: ## Run tests using ginkgo. make dev should have been called ahead of this.
	@ORIG_URI=$$(kubectl -n vetchidev get secret postgres-app -o jsonpath='{.data.uri}' | base64 -d); \
	MOD_URI=$$(echo $$ORIG_URI | sed 's/postgres-rw.vetchidev/localhost/g'); \
	POSTGRES_URI=$$MOD_URI ginkgo -v ./dolores/...

seed: ## Seed the development database
	@ORIG_URI=$$(kubectl -n vetchidev get secret postgres-app -o jsonpath='{.data.uri}' | base64 -d); \
	MOD_URI=$$(echo $$ORIG_URI | sed 's/postgres-rw.vetchidev/localhost/g'); \
	cd dev-seed && POSTGRES_URI=$$MOD_URI go run .

lib: ## Build TypeSpec and install dependencies
	cd typespec && tsp compile . && npm run build && \
	cd ../harrypotter && npm install ../typespec && \
	cd ../ronweasly && npm install ../typespec

docker: ## Build local Docker images for a single platform where it is run
	docker buildx build --load -f harrypotter/Dockerfile-optimized \
		-t psankar/vetchi-harrypotter:$(GIT_SHA) \
		--build-arg API_ENDPOINT="http://hermione:8080" .
	docker buildx build --load -f ronweasly/Dockerfile-optimized \
		-t psankar/vetchi-ronweasly:$(GIT_SHA) \
		--build-arg API_ENDPOINT="http://hermione:8080" .
	docker buildx build --load -f api/Dockerfile-hermione \
		-t psankar/vetchi-hermione:$(GIT_SHA) .
	docker buildx build --load -f api/Dockerfile-granger \
		-t psankar/vetchi-granger:$(GIT_SHA) .
	docker buildx build --load -f sqitch/Dockerfile \
		-t psankar/vetchi-sqitch:$(GIT_SHA) sqitch

publish: ## Build multi-platform Docker images and publish them to the container registry
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "Error: There are uncommitted changes. Please commit them before publishing docker images."; \
		exit 1; \
	fi
	docker buildx inspect multi-platform-builder >/dev/null 2>&1 || docker buildx create --name multi-platform-builder --platform=linux/amd64,linux/arm64 --use
	docker buildx build -f harrypotter/Dockerfile-optimized \
		-t psankar/vetchi-harrypotter:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--build-arg API_ENDPOINT="http://hermione:8080" \
		--push .
	docker buildx build -f ronweasly/Dockerfile-optimized \
		-t psankar/vetchi-ronweasly:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--build-arg API_ENDPOINT="http://hermione:8080" \
		--push .
	docker buildx build -f api/Dockerfile-hermione \
		-t psankar/vetchi-hermione:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--push .
	docker buildx build -f api/Dockerfile-granger \
		-t psankar/vetchi-granger:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--push .
	docker buildx build -f sqitch/Dockerfile \
		-t psankar/vetchi-sqitch:$(GIT_SHA) \
		--platform=linux/amd64,linux/arm64 \
		--push sqitch

devtest: docker ## Brings up an environment with the local docker images. No live reload.
	kubectl delete namespace vetchidevtest --ignore-not-found --force --grace-period=0
	kubectl create namespace vetchidevtest
	kubectl apply --server-side --force-conflicts -f devtest-env/cnpg-1.25.1.yaml
	echo "Waiting for CNPG operator to be ready..."
	sleep 20 && kubectl wait --for=condition=Available deployment/cnpg-controller-manager -n cnpg-system --timeout=5m

	# Then apply core infrastructure
	kubectl apply -n vetchidevtest -f devtest-env/full-access-cluster-role.yaml
	kubectl apply -n vetchidevtest -f devtest-env/postgres-cluster.yaml
	kubectl apply -n vetchidevtest -f devtest-env/minio.yaml
	kubectl apply -n vetchidevtest -f devtest-env/mailpit.yaml
	kubectl apply -n vetchidevtest -f devtest-env/secrets.yaml

	sleep 20 && kubectl wait --for=condition=Ready pod/postgres-1 -n vetchidevtest --timeout=5m
	kubectl wait --for=condition=Ready pod -l app=minio -n vetchidevtest --timeout=5m
	kubectl wait --for=condition=Ready pod -l app=mailpit -n vetchidevtest --timeout=5m

	GIT_SHA=$(GIT_SHA) NAMESPACE=vetchidevtest envsubst '$$GIT_SHA $$NAMESPACE' < devtest-env/sqitch.yaml | kubectl apply -n vetchidevtest -f -
	echo "Waiting for sqitch job to complete..."
	kubectl wait --for=condition=complete job -l app=sqitch -n vetchidevtest --timeout=5m

	# Then apply backend services
	GIT_SHA=$(GIT_SHA) NAMESPACE=vetchidevtest envsubst '$$GIT_SHA $$NAMESPACE' < devtest-env/granger.yaml | kubectl apply -n vetchidevtest -f -
	GIT_SHA=$(GIT_SHA) NAMESPACE=vetchidevtest envsubst '$$GIT_SHA $$NAMESPACE' < devtest-env/hermione.yaml | kubectl apply -n vetchidevtest -f -
	# Finally apply frontend services
	GIT_SHA=$(GIT_SHA) NAMESPACE=vetchidevtest envsubst '$$GIT_SHA $$NAMESPACE' < devtest-env/harrypotter.yaml | kubectl apply -n vetchidevtest -f -
	GIT_SHA=$(GIT_SHA) NAMESPACE=vetchidevtest envsubst '$$GIT_SHA $$NAMESPACE' < devtest-env/ronweasly.yaml | kubectl apply -n vetchidevtest -f -

staging-init: ## Initialize staging environment infrastructure
	kubectl create namespace vetchistaging
	kubectl apply --server-side --force-conflicts -f staging-env/cnpg-1.25.1.yaml
	echo "Waiting for CNPG operator to be ready..."
	kubectl wait --for=condition=Available deployment/cnpg-controller-manager -n cnpg-system --timeout=5m

	# Then apply core infrastructure
	kubectl apply -f staging-env/full-access-cluster-role.yaml
	kubectl apply -f staging-env/postgres-cluster.yaml
	kubectl apply -f staging-env/secrets.yaml

	kubectl get pods -A
	echo "Waiting for postgres to be ready..."
	sleep 10 && kubectl wait --for=condition=Ready pod/postgres-1 -n vetchistaging --timeout=5m

staging: ## Deploy to staging environment
	publish
	GIT_SHA=$(GIT_SHA) envsubst '$$GIT_SHA' < staging-env/sqitch.yaml | kubectl apply -f -
	echo "Waiting for sqitch job to complete..."
	kubectl wait --for=condition=complete job -l app=sqitch -n vetchistaging --timeout=5m

	# Then apply backend services
	GIT_SHA=$(GIT_SHA) envsubst '$$GIT_SHA' < staging-env/granger.yaml | kubectl apply -f -
	GIT_SHA=$(GIT_SHA) envsubst '$$GIT_SHA' < staging-env/hermione.yaml | kubectl apply -f -
	# Finally apply frontend services
	GIT_SHA=$(GIT_SHA) envsubst '$$GIT_SHA' < staging-env/harrypotter.yaml | kubectl apply -f -
	GIT_SHA=$(GIT_SHA) envsubst '$$GIT_SHA' < staging-env/ronweasly.yaml | kubectl apply -f -
