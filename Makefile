obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver receiver/*.go
	@./bin/receiver

calculator:
	@go build -o bin/calculator distance_calculator/*.go
	@./bin/calculator

.PHONY:	obu receiver calculator