gen:
	@jet -dsn="postgresql://bot:bot@localhost:5432/bot?sslmode=disable" -schema=public -path=./internal/repository/gen