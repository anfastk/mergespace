DB_URL=postgres://postgres:131020@localhost:5432/profile_db?sslmode=disable

# make migrate_create name=TABLE_NAME
migrate_create:
	migrate create -ext sql -dir migrations -seq $(name)

# make migrate_up
migrate_up:
	migrate -path migrations -database "$(DB_URL)" up

# make migrate_down
migrate_down:
	migrate -path migrations -database "$(DB_URL)" down