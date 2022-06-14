wasm:
	GOOS=js GOARCH=wasm go build -o demo/public/main.wasm ./sdk-wasm/wasm.go
server:
	go build -o build/server ./server/server.go
setup-db:
	sudo -u postgres bash -c "psql -h 127.0.0.1 -p 5432 < ./server/setup.sql"
	sudo -u postgres bash -c "psql -h 127.0.0.1 -p 5432 -d vault < ./server/vault.sql"