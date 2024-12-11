all:
	go run cmd/main/main.go

migrate-up:
	go run cmd/migrator/main.go --up=true --down=false

migrate-down:
	go run cmd/migrator/main.go --down=true --up=false
