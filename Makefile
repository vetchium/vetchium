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

docker:
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "Error: There are uncommitted changes. Please commit them before building docker images."; \
		exit 1; \
	fi
	$(eval GIT_SHA=$(shell git rev-parse --short=18 HEAD))
	docker build -f Dockerfile-harrypotter -t vetchi/harrypotter:$(GIT_SHA) .
	docker build -f Dockerfile-ronweasly -t vetchi/ronweasly:$(GIT_SHA) .
	docker build -f api/Dockerfile-hermione -t vetchi/hermione:$(GIT_SHA) .
	docker build -f api/Dockerfile-granger -t vetchi/granger:$(GIT_SHA) .
	docker build -f sqitch/Dockerfile -t vetchi/sqitch:$(GIT_SHA) sqitch
	docker tag vetchi/harrypotter:$(GIT_SHA) vetchi/harrypotter:latest
	docker tag vetchi/ronweasly:$(GIT_SHA) vetchi/ronweasly:latest
	docker tag vetchi/hermione:$(GIT_SHA) vetchi/hermione:latest
	docker tag vetchi/granger:$(GIT_SHA) vetchi/granger:latest
	docker tag vetchi/sqitch:$(GIT_SHA) vetchi/sqitch:latest

devtest: docker
	@$(eval GIT_SHA=$(shell git rev-parse --short=18 HEAD))
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
