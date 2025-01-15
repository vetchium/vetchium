module github.com/psankar/vetchi/dev-seed

go 1.23.2

require (
	github.com/jackc/pgx/v5 v5.5.1
	github.com/psankar/vetchi/typespec v0.0.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/text v0.9.0 // indirect
)

replace github.com/psankar/vetchi/typespec => ../typespec
