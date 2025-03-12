$(eval GIT_SHA=$(shell git rev-parse --short=18 HEAD))

dev:
	tilt down
	time kubectl delete pvc -n vetchidev --all --ignore-not-found
	kubectl delete pv -n vetchidev --all --ignore-not-found
	kubectl delete namespace vetchidev --ignore-not-found --force --grace-period=0
	kubectl create namespace vetchidev
	tilt up

test:
	@ORIG_URI=$$(kubectl -n vetchidev get secret postgres-app -o jsonpath='{.data.uri}' | base64 -d); \
	MOD_URI=$$(echo $$ORIG_URI | sed 's/postgres-rw.vetchidev/localhost/g'); \
	POSTGRES_URI=$$MOD_URI ginkgo -v ./dolores/...

seed:
	@ORIG_URI=$$(kubectl -n vetchidev get secret postgres-app -o jsonpath='{.data.uri}' | base64 -d); \
	MOD_URI=$$(echo $$ORIG_URI | sed 's/postgres-rw.vetchidev/localhost/g'); \
	cd dev-seed && POSTGRES_URI=$$MOD_URI go run .

lib:
	cd typespec && tsp compile . && npm run build && \
	cd ../harrypotter && npm install ../typespec && \
	cd ../ronweasly && npm install ../typespec

# Build local images for the host platform
docker:
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "Error: There are uncommitted changes. Please commit them before building docker images."; \
		exit 1; \
	fi
	docker buildx build --load -f Dockerfile-harrypotter -t psankar/vetchi-harrypotter:$(GIT_SHA) .
	docker buildx build --load -f Dockerfile-ronweasly -t psankar/vetchi-ronweasly:$(GIT_SHA) .
	docker buildx build --load -f api/Dockerfile-hermione -t psankar/vetchi-hermione:$(GIT_SHA) .
	docker buildx build --load -f api/Dockerfile-granger -t psankar/vetchi-granger:$(GIT_SHA) .
	docker buildx build --load -f sqitch/Dockerfile -t psankar/vetchi-sqitch:$(GIT_SHA) sqitch

# Build multi-platform images and push them to registry
publish:
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "Error: There are uncommitted changes. Please commit them before publishing docker images."; \
		exit 1; \
	fi
	docker buildx inspect multi-platform-builder >/dev/null 2>&1 || docker buildx create --name multi-platform-builder --platform=linux/amd64,linux/arm64 --use
	docker buildx build --platform=linux/amd64,linux/arm64 -f Dockerfile-harrypotter -t psankar/vetchi-harrypotter:$(GIT_SHA) --push .
	docker buildx build --platform=linux/amd64,linux/arm64 -f Dockerfile-ronweasly -t psankar/vetchi-ronweasly:$(GIT_SHA) --push .
	docker buildx build --platform=linux/amd64,linux/arm64 -f api/Dockerfile-hermione -t psankar/vetchi-hermione:$(GIT_SHA) --push .
	docker buildx build --platform=linux/amd64,linux/arm64 -f api/Dockerfile-granger -t psankar/vetchi-granger:$(GIT_SHA) --push .
	docker buildx build --platform=linux/amd64,linux/arm64 -f sqitch/Dockerfile -t psankar/vetchi-sqitch:$(GIT_SHA) --push sqitch

devtest: docker
	kubectl delete namespace vetchidevtest --ignore-not-found --force --grace-period=0
	kubectl create namespace vetchidevtest
	kubectl apply --server-side --force-conflicts -f devtest-env/cnpg-1.25.1.yaml
	echo "Waiting for CNPG operator to be ready..."
	kubectl wait --for=condition=Available deployment/cnpg-controller-manager -n cnpg-system --timeout=5m

	# Then apply core infrastructure
	kubectl apply -f devtest-env/full-access-cluster-role.yaml
	kubectl apply -f devtest-env/postgres-cluster.yaml
	kubectl apply -f devtest-env/minio.yaml
	kubectl apply -f devtest-env/mailpit.yaml
	kubectl apply -f devtest-env/secrets.yaml

	sleep 5 && kubectl wait --for=condition=Ready pod/postgres-1 -n vetchidevtest --timeout=5m
	kubectl wait --for=condition=Ready pod -l app=minio -n vetchidevtest --timeout=5m
	kubectl wait --for=condition=Ready pod -l app=mailpit -n vetchidevtest --timeout=5m

	GIT_SHA=$(GIT_SHA) envsubst '$$GIT_SHA' < devtest-env/sqitch.yaml | kubectl apply -f -
	echo "Waiting for sqitch job to complete..."
	kubectl wait --for=condition=complete job/sqitch -n vetchidevtest --timeout=5m

	# Then apply backend services
	GIT_SHA=$(GIT_SHA) envsubst '$$GIT_SHA' < devtest-env/granger.yaml | kubectl apply -f -
	GIT_SHA=$(GIT_SHA) envsubst '$$GIT_SHA' < devtest-env/hermione.yaml | kubectl apply -f -
	# Finally apply frontend services
	GIT_SHA=$(GIT_SHA) envsubst '$$GIT_SHA' < devtest-env/harrypotter.yaml | kubectl apply -f -
	GIT_SHA=$(GIT_SHA) envsubst '$$GIT_SHA' < devtest-env/ronweasly.yaml | kubectl apply -f -

staging-init:
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

staging: publish
	GIT_SHA=$(GIT_SHA) envsubst '$$GIT_SHA' < staging-env/sqitch.yaml | kubectl apply -f -
	echo "Waiting for sqitch job to complete..."
	kubectl wait --for=condition=complete job -l app=sqitch -n vetchistaging --timeout=5m

	# Then apply backend services
	GIT_SHA=$(GIT_SHA) envsubst '$$GIT_SHA' < staging-env/granger.yaml | kubectl apply -f -
	GIT_SHA=$(GIT_SHA) envsubst '$$GIT_SHA' < staging-env/hermione.yaml | kubectl apply -f -
	# Finally apply frontend services
	GIT_SHA=$(GIT_SHA) envsubst '$$GIT_SHA' < staging-env/harrypotter.yaml | kubectl apply -f -
	GIT_SHA=$(GIT_SHA) envsubst '$$GIT_SHA' < staging-env/ronweasly.yaml | kubectl apply -f -
