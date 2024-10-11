# vetchi

## Development

### Prerequisites

- [Tilt](https://docs.tilt.dev/install.html)
- [Docker](https://docs.docker.com/get-docker/)
- [Kubernetes](https://kubernetes.io/docs/tasks/tools/)

### Setup

To bring the services up, run the following commands:
```
$ kubectl create namespace vetchidev
$ tilt up
```

To connect to the port-forwarded Postgres using the psql command line, use the following command:
```
$ psql -h localhost -p 5432 -U user vdb
```

### Tear down

To tear down the services, run the following command:
```
$ tilt down
$ kubectl delete namespace vetchidev
```

### Notes

- [hermione](api/hermione) contains the stateless API server that can be scaled horizontally. Almost all HTTP handlers should be implemented here.
- [granger](api/granger) contains the singleton API server with global variables, that should NOT be scaled horizontally. Almost no HTTP handler should be implemented here. This should be used for periodic tasks and other such bookkeeping.
- [hermione](api/hermione) and [granger](api/granger) share the same go.mod and go.sum and together they implement the Vetchi API
- [harrypotter](harrypotter) contains the React.js frontend for the Employer app
- [ronweasly](ronweasly) contains the React.js frontend for the Hub app
- [sqitch](sqitch) contains the database migration scripts
- [dolores](dolores) contains the end to end tests for the API server
- Use [Ginkgo](https://onsi.github.io/ginkgo/) for writing tests
- Use [Gomega](https://onsi.github.io/gomega/) for assertions
- Use [goline](https://github.com/segmentio/golines) to format the Go code. Do NOT manually format the code or split the parameters to multiple lines.
- Use [prettier](https://prettier.io/) to format the typescript code. Do not manually format the code or split the parameters to multiple lines.
- Use the below snippet to sort the openapi-spec.yaml file
```
$ yq eval 'sort_keys(..)' vetchi-openapi.yml -o=yaml > output.yaml
$ # Move the openapi and info tags to the top of the file
$ # Ensure the markdown is valid by editor.swagger.io or editor plugins
$ mv output.yaml vetchi-openapi.yml
$ # Alternatively you can use a custom yaml sort in your editor
```
