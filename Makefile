migrations:
	goose -dir db/migrations postgres "postgresql://jir:root@127.0.0.1:5432/projects_service?sslmode=disable" up
dropmigrations:
	goose -dir db/migrations postgres "postgresql://jir:root@127.0.0.1:5432/projects_service?sslmode=disable" down
run:
	go build -o projects_service && ./projects_service