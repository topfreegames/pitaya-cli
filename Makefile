build:
	@mkdir -p out
	@go build -o ./out/pitaya-cli-darwin ./main.go

build-linux:
	@mkdir -p out
	@GOOS=linux GOARCH=amd64 go build -o ./out/pitaya-cli-linux ./main.go
