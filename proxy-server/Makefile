default: run

PROXY_ADDR ?="localhost:5433"
DATABASE_ADDR ?="localhost:5432"
STORAGE_PATH ?="/Users/egorskurihin/desktop/Git/postgres-proxy/sessions.json"

run:
	go run main.go \
		-db-addr=$(DATABASE_ADDR) \
		-proxy-addr=$(PROXY_ADDR) \
		-storage-path=$(STORAGE_PATH)
