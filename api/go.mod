module github.com/psankar/vetchi/api

go 1.23.2

require (
	github.com/fatih/color v1.18.0
	github.com/go-playground/validator/v10 v10.22.1
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.1
	github.com/psankar/vetchi/typespec v0.0.0
	github.com/wneessen/go-mail v0.5.0
	golang.org/x/crypto v0.28.0
)

replace github.com/psankar/vetchi/typespec => ../typespec

require (
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
)
