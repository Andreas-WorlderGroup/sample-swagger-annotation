run-container-core:
	docker rm core-container --force
	docker build -t core .
	docker run -d -it -p 8080:8080 --name=core-container core

run-worker:
	go run cmd/worker/main.go

generate:
	swag init -g cmd/core/main.go -o ./docs

install:
	go get \
		github.com/swaggo/swag/cmd/swag \
		github.com/swaggo/echo-swagger \
		github.com/go-co-op/gocron