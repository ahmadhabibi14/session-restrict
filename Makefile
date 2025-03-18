setup:
	go install github.com/air-verse/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	sudo curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64
	sudo chmod +x /usr/local/bin/dbmate

test:
	go test ./tests/... -v

test-coverage:
	go test -coverpkg=./src/repo/... ./tests