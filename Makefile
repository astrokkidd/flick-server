MESSAGE = migration

## help: print this help message
.PHONY: help
help:
	@sed -n "s/^## \(.*\): \(.*\)/\1: \2/p" ${MAKEFILE_LIST}

## clean: clean the build files
.PHONY: clean
clean:
	rm -rf ./tmp
	rm ./pkg/database/*.sql.go


## generate: generate SQL code
.PHONY: generate
generate:
	go generate -v ./...

## build: build the Go application
.PHONY: build
build: generate
	go build -v -o ./tmp ./...

## migrate/diff: create new migration from schema diffs
.PHONY: migrate/diff
migrate/diff:
	atlas migrate diff $(MESSAGE) --dir "file://migrations" --to "file://configs/schema/" --dev-url "docker://postgres/17/dev"

## migrate/apply: apply new migration
.PHONY: migrate/apply
migrate/apply:
	atlas migrate apply --url "postgres://postgres:admin@localhost:5432/flick?sslmode=disable"   

## watch: launch with live reload
.PHONY: watch
watch:
	air -c ./configs/.air.toml