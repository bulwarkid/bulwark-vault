wasm:
	GOOS=js GOARCH=wasm go build -o demo/public/main.wasm ./sdk-wasm/wasm.go
build/server:
	go build -o build/server ./server/server.go
setup-db:
	sudo -u postgres bash -c "psql < ./server/setup.sql"
	sudo -u postgres bash -c "psql -d vault < ./server/vault.sql"