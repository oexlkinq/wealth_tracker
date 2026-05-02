.PHONY: setup build migrate clean sql

# TODO: добавить запуск sqlc перед сборкой
setup: clean migrate build
	./wealth_tracker setup example > data.yaml
	./wealth_tracker setup data.yaml
	rm -f data.yaml

build:
	go build
