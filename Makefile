test:
	go test ./... -timeout 15s -race -cover -coverprofile=coverage.out -v
	go tool cover -html=coverage.out -o coverage.html

pg:
	docker pull postgres:13-alpine
	docker run \
		-d \
		--name test.pg \
		-e POSTGRES_PASSWORD=596a96cc7bf9108cd896f33c44aedc8a \
		-e POSTGRES_DB=testpg \
		-e POSTGRES_USER=postgres \
		-p 7005:5432 \
		postgres:13-alpine

pg_rm:
	docker rm -f test.pg
