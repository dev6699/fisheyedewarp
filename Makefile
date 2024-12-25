.PHONY: run
run:
	go run --tags opencvstatic cmd/fisheyedewarp/main.go

.PHONY: build
build:
	go build --tags opencvstatic -o fisheyedewarp cmd/fisheyedewarp/main.go