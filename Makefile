PROJECT = aws-rotate-key

build:
	go fmt
	go build -o $(GOPATH)/bin/$(PROJECT) github.com/scottbrown/$(PROJECT)

test:
	$(PROJECT)

