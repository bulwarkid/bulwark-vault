wasm:
	GOOS=js GOARCH=wasm go build -o demo/public/main.wasm ./sdk-wasm/wasm.go
build:
	if [ -d output ]; then rm output/*; fi
	if [ ! -d output]; then mkdir output; fi
	go build -o output/server ./server/server.go
	npm run build --prefix demo
setup-db:
	sudo -u postgres bash -c "psql < ./server/setup.sql"
	sudo -u postgres bash -c "psql -d vault < ./server/vault.sql"