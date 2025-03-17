secret_settings(disable_scrub=True)

k8s_kind('Cluster', api_version='postgresql.cnpg.io/v1')

# Load Kubernetes YAML files
k8s_yaml('tilt-env/full-access-cluster-role.yaml')
k8s_yaml('tilt-env/postgres-cluster.yaml')
k8s_yaml('tilt-env/sqitch.yaml')
k8s_yaml('tilt-env/mailpit.yaml')
k8s_yaml('tilt-env/minio.yaml')
k8s_yaml('tilt-env/secrets.yaml')
k8s_yaml('tilt-env/hermione.yaml')
k8s_yaml('tilt-env/granger.yaml')
k8s_yaml('tilt-env/harrypotter.yaml')
k8s_yaml('tilt-env/ronweasly.yaml')
k8s_yaml('tilt-env/sortinghat.yaml')

# Define Docker builds with root context to include typespec
docker_build('psankar/vetchi-granger', '.', dockerfile='api/Dockerfile-granger')
docker_build('psankar/vetchi-hermione', '.', dockerfile='api/Dockerfile-hermione')
docker_build('psankar/vetchi-sqitch', 'sqitch', dockerfile='sqitch/Dockerfile')
docker_build('psankar/vetchi-sortinghat', '.', dockerfile='sortinghat/Dockerfile')

# Development builds for Next.js apps with live reload
docker_build(
    'psankar/vetchi-harrypotter',
    '.',
    dockerfile='harrypotter/Dockerfile-tilt',
    build_args={'API_ENDPOINT': 'http://hermione:8080'}
)

docker_build(
    'psankar/vetchi-ronweasly',
    '.',
    dockerfile='ronweasly/Dockerfile-tilt',
    build_args={'API_ENDPOINT': 'http://hermione:8080'}
)

# Infra services - postgres is port-forwarded differently because of the
# cnpg operator packaging semantics with multiple deployments, pods coming up.
k8s_resource('mailpit', port_forwards='8025:8025')
k8s_resource('minio', port_forwards='9000:9000')

# Backend API services
k8s_resource('hermione', port_forwards='8080:8080')
k8s_resource('granger', port_forwards='8081:8080')
k8s_resource('sortinghat', port_forwards='8082:8080')

# Frontend apps
k8s_resource('harrypotter', port_forwards='3001:3000')
k8s_resource('ronweasly', port_forwards='3002:3000')


# The cnpg operator takes a lot of time for the pg pods to get ready
# So we need to do all the below magic for pg port_forward alone unlike
# the rest of the services done above.
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
    deps=["devtest-env/postgres-cluster.yaml"],  # Ensures Postgres exists before starting port-forward
    allow_parallel=True,
    serve_cmd="kubectl -n vetchidev port-forward service/postgres-rw 5432:5432"
)