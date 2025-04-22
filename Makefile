build:
	@go build -o bin/ecom cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/ecom

GOMIGRATE := $(shell go env GOPATH)/bin/migrate
DB_URL := postgres://test:test@localhost:5432/test?sslmode=disable

migration:
	@$(GOMIGRATE) create -ext sql -dir migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	$(GOMIGRATE) -path ./migrations -database "$(DB_URL)" up

migrate-down:
	$(GOMIGRATE) -path ./migrations -database "$(DB_URL)" down 1

migrate-reset:
	$(GOMIGRATE) -path ./migrations -database "$(DB_URL)" down

migrate-force:
	$(GOMIGRATE) -path ./migrations -database "$(DB_URL)" force $(version)

migrate-drop:
	$(GOMIGRATE) -path ./migrations -database "$(DB_URL)" drop -f

%:
	@:
