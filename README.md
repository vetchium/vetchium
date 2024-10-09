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
$ psql -h localhost -p 5432 -U user vetchidb
```

### Tear down

To tear down the services, run the following command:
```
$ tilt down
$ kubectl delete namespace vetchidev
```

### Notes

- [hermione](api/hermione) contains the stateless API server that can be scaled horizontally
- [granger](api/granger) contains the singleton API server with global variables, that should NOT be scaled horizontally
- [hermione](api/hermione) and [granger](api/granger) share the same go.mod and go.sum and together they implement the Vetchi API
- [harrypotter](harrypotter) contains the React.js frontend for the Employer app
- [ronweasly](ronweasly) contains the React.js frontend for the Hub app
- [sqitch](sqitch) contains the database migration scripts
