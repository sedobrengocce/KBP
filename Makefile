build:
	@go build -o ./dist/kbp -v ./cmd 
run:
	@go run -v ./cmd	
install:
	@cp ./dist/kbp /usr/bin/
