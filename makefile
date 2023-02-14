export DEVELOPMENT_MODE=true

test: export CONFIG_PATH=..
test-r: export CONFIG_PATH=..

test:
	go test -count=1 -v ./tests

# make test-r run=TestApiMsg
test-r:
	go test -count=1 -v ./tests -run "$(run)"

run:
	go run main.go

run-b:
	bin/app

.PHONY: docker-up-d
docker-up-d:
	docker-compose -f docker-compose.yaml up --build -d

.PHONY: build
build:
	go build -o bin/app main.go