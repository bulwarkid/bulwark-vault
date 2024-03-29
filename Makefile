wasm:
	GOOS=js GOARCH=wasm go build -o demo/public/main.wasm ./sdk-wasm/wasm.go
build:
	./scripts/build.sh
setup-db:
	sudo -u postgres bash -c "psql < ./server/setup.sql"
	sudo -u postgres bash -c "psql -d vault < ./server/vault.sql"
format:
	gofmt -l -w ./sdk/*.go
	gofmt -l -w ./server/*.go
	gofmt -l -w ./sdk-wasm/*.go