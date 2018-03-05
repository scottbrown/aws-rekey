PROJECT = aws-rotate-key

build:
	go fmt
	go build -o $(PROJECT) github.com/scottbrown/$(PROJECT)

test:
	./$(PROJECT)

