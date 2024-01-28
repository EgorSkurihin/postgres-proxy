default: run

ADDR ?=":80"
STORAGE_PATH ?="/Users/egorskurihin/Desktop/Git/postgres-proxy/sessions.json"
USERS_FILE_PATH ?="/Users/egorskurihin/Desktop/Git/postgres-proxy/web-server/users.json"

run:
	go run main.go \
		-addr=$(ADDR) \
		-storage-path=$(STORAGE_PATH) \
		-users-file-path=$(USERS_FILE_PATH)
