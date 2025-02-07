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
vetchi $ # Bring up the backend services 
vetchi $ make
vetchi $ # Visit http://localhost:10350/ to see the tilt UI which will show you the services, logs, port-forwards, etc.

vetchi $ cd typespec && npm ci && npm run build

vetchi $ # Bring up the frontend services
vetchi $ cd harrypotter && npm ci && npm run dev
vetchi $ cd ronweasly && npm ci && npm run dev
vetchi $ # Visit http://localhost:3000 and http://localhost:3001

vetchi $ # If changes were made in typespec/**/*.ts files, do:
vetchi $ cd typespec && npm run build 
vetchi $ cd harrypotter && make ;# This installs the new deps from typespec and does `npm run dev`
vetchi $ cd ronweasly && make

vetchi $ # Seed some test data
vetchi $ make seed  # tilt up should be running
```

```
http://localhost:3000 contains the Employer site. Login with:
domain: gryffindor.example
username: hermione@gryffindor.example (or) admin@gryffindor.example
password: NewPassword123$
```

```
http://localhost:3001 contains the Hub site. Login with:
user: hagrid@hub.example (or) minerva@hub.example
password: NewPassword123$
```

* Any Openings created on the Employer site should be available on the Hub site for Users to Find Opening and Apply
* Any Applications made on Hub site should be available for the Employer to either Shortlist or Reject
* Any shortlisted application will become a Candidacy which both the Employer and Hub User can use for further communication. The url for the candidacy will be in the mail at http://localhost:8025 

To connect to the port-forwarded Postgres using psql, get the connection details from the Kubernetes secret:
```
$ POSTGRES_URI=$(kubectl -n vetchidev get secret postgres-app -o jsonpath='{.data.uri}' | base64 -d | sed 's/postgres-rw.vetchidev/localhost/g')
$ psql "$POSTGRES_URI"
```

To connect to the port-forwarded Postgres using DBeaver or some such JDBC client, use:
```
$ kubectl -n vetchidev get secret postgres-app -o jsonpath='{.data.uri}' | base64 -d | sed -E 's|postgresql://([^:]+):([^@]+)@([^:/]+):([0-9]+)/([^?]+)|jdbc:postgresql://localhost:\4/\5?user=\1\&password=\2|'
```

To run tests, use the following command:
```
$ go install github.com/onsi/ginkgo/v2/ginkgo; # Only once
vetchi $ make test ; # tilt up should be running
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
- [hedwig](api/internal/hedwig) is a library for sending emails. There is a template folder that contains the templates for the emails. Each template has a name, a html file and a text file. The template values should be consistent between the html and text files. create-onboard-emails does not use hedwig yet and should be migrated.
- Use [Ginkgo](https://onsi.github.io/ginkgo/) for writing API tests
- Use [Gomega](https://onsi.github.io/gomega/) for assertions
- Use [golines](https://github.com/segmentio/golines) to format the Go code. Do NOT manually format the code or split the parameters to multiple lines. Write a really long line with all parameters and then summon golines to format it.
- Use [prettier](https://prettier.io/) to format the typescript code. Do not manually format the code or split the parameters to multiple lines.
- [dev-seed](dev-seed) contains a sample set of employers, hub users, openings, etc. Feel free to extend this to cover more scenarios. The initial employer data is hard-coded directly to the database. The rest are all created via APIs. So you need `tilt up` to be running before this. This will also serve as some basic sanity testing.

## Engineering Notes
Following are some rules that you should follow while working on the code. It is okay to break these rules, if that would make the code more readable. But your interest to break rules should not spawn from your inability to follow rules.

#### Backend
- Readability > Scalability > Performance. Optimise in this order.
- We are a simple CRUD application that just needs to scale well. No need for anything fancy. Use boring, stable, battle-tested technologies. Try to keep things simple but do not make anything that would be impossible to scale later without a lot of rework. An example of this is: We chose postgres and neither sqlite (difficult to scale later for multiple machines) nor spanner (huh !?)
- Do not use fancy algorithms. Use simple and scalable solutions.
- Do not use ORMs. Do not fear SQL.
- Do not introduce a caching layer (like Redis or within the Golang handlers) for any data that has a strong consistency requirement. Rely on database indexes and query tuning for performance. It is okay to put a lot of load on the database. It is easier to scale a single ACID compliant database layer (either vertically or horizontally) than debug a problem across multiple caching layers.
- Group related items in Code together. For example a section could be for Location related functions, another section could be for OrgUser related functions. Within a group, try to maintain some kind of pattern that would be obvious (like alphabetical ordering of methods) and consistent across multiple groups/files.
- Do not depend on any library unnecessarily (only one exception mentioned below). Try to reimplement in simple Go
- Do not reimplement any security related features/algorithms yourself. Use well-established libraries and algorithms. Eg: Use bcrypt not your own hashing algorithm.
- Do not create more modules for the backend. Try to code within one of Hermione or Granger.
- Do not use any kubernetes specific abstractions. Eg: Do not create a Kubernetes Job to send email, but instead use goroutines and channels. The lesser the images that we have, the easier it is to maintain.
- All configuration data should be read from kubernetes configmap
- All sensitive data (passwords, API keys, etc.) should be read from kubernetes secrets
- All backend APIs should have test coverage. Write exhaustive tests for positive and negative cases, border conditions. Focus on meaningfully detecting regressions and not just on coverage percentages.
- End to end tests > Unit tests. It is okay to ignore unit tests as long as the end to end tests are comprehensive across all codepaths.
- We use [ginkgo](https://onsi.github.io/ginkgo/) for writing end to end tests. Each test should have a testcase-up.sql and testcase-down.sql file. The testcase-up.sql file should be used to setup the test data and the testcase-down.sql should be used to clean up the test data. All testcases must clean whatever data they create (including emails). All testcases must be [idempotent](https://en.wikipedia.org/wiki/Idempotence#Computer_science_meaning).
- Minimize data that has to be moved out of database to backend. But have most business logic in Go code. This may seem contradictory at first, but if you read through the code, you will understand.
- Use the vetchi.Structs to pass data around everywhere as much as possible. If a handler has to write to multiple tables for a single HTTP request, then create and use a new struct under the `db` package. All things under the `db` package should be strictly internal and should not be exposed to the API (for frontend or clients). Things under the `db` package can be changed anytime.
- We use [typespec](https://typespec.io/) to define our API contract. Backend and Frontend code should adhere to the typespec specification. Once the typespec is compiled, it will export an openapi.yaml file, which can be previewed in your favorite editor. Take a quick 5 minute introduction to typespec, if you have not worked with it before. It is quite easy to pick up. It is more concise than yaml and is far easier to edit/read.
- Write the typespec specification first before writing any new code. It is okay to change the spec until the code is merged, but should be considered set in stone after that.
- End all files with a newline. Do NOT have any trailing whitespace.
- Enforce best-practices via editorconfig, CI or other FOSS tooling automation as much as possible. It is the responsibility of the reviewers to check for these.
- Merge small changes frequently. Hide things behind feature flags until they are tested for functionality and scale. Do not drop big changes.
- Maintain 80 column limits for all the code. The test files under dolores do not have a 80 column limit. Sometimes the SQL under the postgres package may also make things difficult to fit under 80 columns which we have to live with. But ensure that the code is readable.
- Try to have a maximum of about 200 lines per file. Be miserly in creating new packages and be generous in creating new files under existing packages.
- Use https://sqlformat.darold.net/ to format SQL within the postgres.go file and dolores/*.pgsql files. This does not do a good job at breaking long line or aligning complex queries, but it is better to be consistently formatted. If there is a better pgsqlfmt in future, use.
- In the backend, log errors with the `Err` method. Log as Error only on the place where the error actually happens. This will help us get maximum debug information. In all the above layers of the call stack, if you want to log, use the `Dbg` method. The only exceptions to this are when the services are coming up. In that case, the errors are logged in the main function. Always, strive to log an error as Err only once. This will help us avoid generating too many tickets in SIEM/EventMgmt systems (such as Sentry).

#### Frontend
- Format all code with prettier
- Do not duplicate any structs to send requests to the backend or parse the responses from the server. Use the library imported from typespec.
- Material UI is used for the frontend styling. Stick to the same styling guidelines to keep the UI consistent. Do not import any new icon families or components from other libraries.
- Some of the libraries may be using deprecated versions. Always try to upgrade to the latest stable releases.


#### Others
Sometimes we are forced to use longer names in the libraries. For example, we use  employer.EmployerInterview instead of just employer.Interview as any Go programmer would do. The reason for this is, typespec does not allow creating duplicate structures as we compile into one large openapi.yaml in the end. So we have to live with the longer names. But try to minimize the length of variables as much as you can.
