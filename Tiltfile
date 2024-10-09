# Load Kubernetes YAML files
k8s_yaml('api/granger-k8s.yaml')
k8s_yaml('api/hermione-k8s.yaml')
k8s_yaml('harrypotter/harrypotter-tilt.yaml')
k8s_yaml('ronweasly/ronweasly-tilt.yaml')

# Define Docker builds
docker_build('granger', 'api', dockerfile='api/Dockerfile-granger')
docker_build('hermione', 'api', dockerfile='api/Dockerfile-hermione')
docker_build('harrypotter', 'harrypotter', dockerfile='harrypotter/Dockerfile')
docker_build('ronweasly', 'ronweasly', dockerfile='ronweasly/Dockerfile')

# Associate images with Kubernetes resources
k8s_resource('granger', port_forwards='8080:8080')
k8s_resource('hermione', port_forwards='8081:8080')
k8s_resource('harrypotter', port_forwards='3000:3000')
k8s_resource('ronweasly', port_forwards='3001:3000')
