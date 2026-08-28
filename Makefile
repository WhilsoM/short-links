.PHONY: run lint test migration-create migration-up migration-down git-log

run:
	go run ./cmd/api


consumer:
	go run ./cmd/analytics

lint:
	golangci-lint run

test:
	go test ./...

migration-create:
	goose -dir migrations create $(name) sql

migration-up:
	goose -dir migrations postgres "$(DATABASE_URL)" up

migration-down:
	goose -dir migrations postgres "$(DATABASE_URL)" down

git-log:
	git log --oneline --decorate --all --graph
