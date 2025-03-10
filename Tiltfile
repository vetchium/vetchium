secret_settings(disable_scrub=True)

k8s_kind('Cluster', api_version='postgresql.cnpg.io/v1')

# Load Kubernetes YAML files
k8s_yaml('tiltenv/cnpg-1.25.1.yaml')
k8s_yaml('tiltenv/postgres-cluster.yaml')
k8s_yaml('tiltenv/mailpit.yaml')
k8s_yaml('tiltenv/minio.yaml')
k8s_yaml('tiltenv/full-access-cluster-role.yaml')
k8s_yaml('tiltenv/secrets.yaml')
k8s_yaml('api/hermione-tilt.yaml')
k8s_yaml('api/granger-tilt.yaml')
k8s_yaml('sqitch/sqitch-tilt.yaml')

# Define Docker builds with root context to include typespec
docker_build('psankar/granger', '.',
    dockerfile='api/Dockerfile-granger',
)

docker_build('psankar/hermione', '.',
    dockerfile='api/Dockerfile-hermione',
)

docker_build('psankar/vetchi-sqitch', 'sqitch', dockerfile='sqitch/Dockerfile')

k8s_resource('mailpit', port_forwards='8025:8025')
k8s_resource('granger', port_forwards='8080:8080')
k8s_resource('hermione', port_forwards='8081:8080')
k8s_resource('sqitch')

# Function to wait for the Kubernetes service to be ready
def wait_for_service(namespace, service):
    print("Waiting for service {service} in namespace {namespace}...")

    while True:
        result = local("kubectl -n {} get service {}".format(namespace, service), quiet=True)
        if "NotFound" not in result and "No resources found" not in result:
            print("Service {service} is now available.")
            break

# Define a local resource to handle port-forwarding
local_resource(
    "postgres-port-forward",
    cmd="sh -c 'while ! kubectl -n vetchidev get service postgres-rw; do sleep 10; done && kubectl -n vetchidev port-forward service/postgres-rw 5432:5432'",
    deps=["tiltenv/postgres-cluster.yaml"],  # Ensures Postgres exists before starting port-forward
    allow_parallel=True,
    serve_cmd="kubectl -n vetchidev port-forward service/postgres-rw 5432:5432"
)