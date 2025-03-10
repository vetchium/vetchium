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
