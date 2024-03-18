run:
	docker-compose up

doc-gen:
	swag init -g ./cmd/serv.go

test:
	go test ./tests -v