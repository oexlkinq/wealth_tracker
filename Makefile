.PHONY: setup build migrate clean sql

# TODO: добавить запуск sqlc перед сборкой
setup: clean migrate build
	./wealth_tracker setup example > data.yaml
	./wealth_tracker setup data.yaml
	rm -f data.yaml

build:
	go build

migrate:
	goose -dir ./internal/db/migrations sqlite3 ./appdata/wealth_tracker.db up

clean:
	rm -rf ./appdata/wealth_tracker.db ./.gen ./wealth_tracker

sql:
	sqlite3 appdata/wealth_tracker.db
