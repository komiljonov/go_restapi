
# Create new migration
migrate create -ext sql -dir db/migrations add_users_table

# Migrate to db
migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/go?sslmode=disable" up
