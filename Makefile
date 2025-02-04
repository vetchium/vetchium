dev:
	tilt down
	@if [ ! -f cnpg-1.24.2.yaml ]; then \
		curl -o cnpg-1.24.2.yaml https://raw.githubusercontent.com/cloudnative-pg/cloudnative-pg/release-1.24/releases/cnpg-1.24.2.yaml; \
	fi
	kubectl delete -f cnpg-1.24.2.yaml --ignore-not-found
	kubectl apply --server-side -f cnpg-1.24.2.yaml
	kubectl delete namespace vetchidev --ignore-not-found
	kubectl create namespace vetchidev
	kubectl wait --for=condition=established --timeout=60s crd/clusters.postgresql.cnpg.io
	kubectl -n cnpg-system wait --for=condition=ready pod -l app.kubernetes.io/name=cloudnative-pg --timeout=60s
	kubectl apply -f postgres-cluster.yaml
	sleep 10
	kubectl -n vetchidev wait --for=condition=ready pod -l cnpg.io/podRole=instance --timeout=240s
	kubectl -n vetchidev port-forward service/postgres-rw 5432:5432 &
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

