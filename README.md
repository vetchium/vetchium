# vetchi

## Development

### Prerequisites

- [Tilt](https://docs.tilt.dev/install.html)
- [Docker](https://docs.docker.com/get-docker/)
- [Kubernetes](https://kubernetes.io/docs/tasks/tools/)
- [Go](https://golang.org/doc/install)

### Setup

To bring the services up, run the following commands:
```
vetchi $ # Setup any Kubernetes cluster (docker desktop, kind, etc.) and make kubectl point to it
vetchi $ kubectl create namespace vetchidev
vetchi $ tilt up
vetchi $ # Visit http://localhost:10350/ to see the tilt UI which will show you the services, logs, port-forwards, etc.
```

To connect to the port-forwarded Postgres using the psql command line, use the following command:
```
$ psql -h localhost -p 5432 -U user vdb
```

To run tests, use the following command:
```
$ go install github.com/onsi/ginkgo/v2/ginkgo; # Only once
vetchi $ cd dolores
dolores $ ginkgo -vv ; # tilt up should be running
```

### Tear down

To tear down the services, run the following command:
```
vetchi $ tilt down
vetchi $ kubectl delete namespace vetchidev
```

### Code Structure
- [hermione](api/hermione) contains the stateless API server that can be scaled horizontally. Almost all HTTP handlers should be implemented here.
- [granger](api/granger) contains the singleton API server with global variables, that should NOT be scaled horizontally. Almost no HTTP handler should be implemented here. This should be used for periodic tasks and other such bookkeeping on the backend.
- [hermione](api/hermione) and [granger](api/granger) share the same go.mod and go.sum and together they implement the Vetchi API
- [harrypotter](harrypotter) contains the React.js frontend for the Employer app
- [ronweasly](ronweasly) contains the React.js frontend for the Hub app
- [sqitch](sqitch) contains the database migration scripts
- [dolores](dolores) contains the end to end tests for the API server
- Use [Ginkgo](https://onsi.github.io/ginkgo/) for writing API tests
- Use [Gomega](https://onsi.github.io/gomega/) for assertions
- Use [golines](https://github.com/segmentio/golines) to format the Go code. Do NOT manually format the code or split the parameters to multiple lines. Write a really long line with all parameters and then summon golines to format it.
- Use [prettier](https://prettier.io/) to format the typescript code. Do not manually format the code or split the parameters to multiple lines.
- Use the below snippet to sort the openapi-spec.yaml file
```
$ yq eval 'sort_keys(..)' vetchi-openapi.yml -o=yaml > output.yaml
$ # Move the openapi and info tags to the top of the file
$ # Ensure the yaml is valid by editor.swagger.io or editor plugins
$ mv output.yaml vetchi-openapi.yml
$ # Alternatively you can use a custom yaml sort in your editor
```

### Engineering Notes
Following are some rules that you should follow while working on the code. It is okay to break these rules, if that would make the code more readable. But your interest to break rules should not stem from your inability to follow rules.

- Readability > Scalability > Performance. Optimise in this order.
- Do not use fancy algorithms. Use simple and scalable solutions.
- Do not use ORMs. Do not fear SQL.
- Do not introduce a caching layer (like Redis). Rely on database indexes and query tuning for performance.
- Always sort the methods in an interface, OpenAPI spec, etc., alphabetically, so that it is easier for editing. Try as much possible to keep any list of items in code alphabetically sorted. There may be exceptions where grouping items together will help with readability. Use your best judgement.
- Do not depend on any library unnecessarily (only one exception mentioned below). Try to reimplement in simple Go or Typescript.
- Do not reimplement any security related features yourself. Use well-established libraries and algorithms. Eg: Use bcrypt not your own hashing algorithm.
- Do not create more modules for the backend. Try to code within one of Hermione or Granger.
- Do not use any kubernetes specific abstractions. Eg: Do not create a Kubernetes Job to send email but use goroutines and channels.
- All configuration data should be read from configmap
- All sensitive data (passwords, API keys, etc.) should be read from secrets
- All backend APIs should have test coverage. Write exhaustive tests for positive and negative cases, border conditions. Focus on meaningfully detecting regressions and not just on coverage percentages.
- End to end tests > Unit tests
- We use [ginkgo](https://onsi.github.io/ginkgo/) for writing end to end tests. Each test should have a testcase-up.sql and testcase-down.sql file. The testcase-up.sql file should be used to setup the test data and the testcase-down.sql should be used to clean up the test data. All testcases must clean whatever data they create (including emails). All testcases must be idempotent.
- Minimize data that has to be moved out of database to backend. But have most business logic in Go code. This may seem contradictory at first, but if you read through the code, you will understand.
- What is mentioned in the openapi spec is the contract. Try to keep the implementation as close to the contract as possible. Backend and Frontend code should adapt to the openapi spec.
- Write openapi spec first before writing any new code. It is okay to change the spec until the code is merged, but should be considered set in stone after that.
- End all files with a newline. Do NOT have any trailing whitespace.
- Enforce best-practices via editorconfig, CI or other FOSS tooling automation as much as possible. It is the responsibility of the reviewers to check for these.
- Merge small changes frequently. Hide things behind feature flags until they are tested for functionality and scale. Do not drop big changes.
- The test files under dolores do not have a 80 column limit. But ensure that the code is readable. Try to maintain 80 column limit for the rest of the code.
