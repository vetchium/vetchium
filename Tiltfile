# Load Kubernetes YAML files
k8s_yaml('vetchi-secrets-tilt.yaml')
k8s_yaml('api/granger-tilt.yaml')
k8s_yaml('api/hermione-tilt.yaml')
k8s_yaml('harrypotter/harrypotter-tilt.yaml')
k8s_yaml('ronweasly/ronweasly-tilt.yaml')
k8s_yaml('sqitch/sqitch-tilt.yaml')
k8s_yaml('sqitch/postgres-tilt.yaml')

# Define Docker builds
docker_build('psankar/granger', 'api', dockerfile='api/Dockerfile-granger')
docker_build('psankar/hermione', 'api', dockerfile='api/Dockerfile-hermione')
docker_build('psankar/harrypotter', 'harrypotter', dockerfile='harrypotter/Dockerfile')
docker_build('psankar/ronweasly', 'ronweasly', dockerfile='ronweasly/Dockerfile')
docker_build('psankar/vetchi-sqitch', 'sqitch', dockerfile='sqitch/Dockerfile')

# Associate images with Kubernetes resources
k8s_resource('granger', port_forwards='8080:8080')
k8s_resource('hermione', port_forwards='8081:8080')
k8s_resource('harrypotter', port_forwards='3000:3000')
k8s_resource('ronweasly', port_forwards='3001:3000')
k8s_resource('sqitch')
k8s_resource('postgres', port_forwards='5432:5432')