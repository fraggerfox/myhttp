#
# A simple Makefile containing basic targets to build and run the program
#
# NOTE: Run make depends to get the testing package.
#

default:
	@make help

build:
	go build -v myhttp.go

clean:
	go clean -x
	rm -f coverage*

test:
	go test -test.v -run ''

cover:
	go test -cover
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

help:
	@echo "------------------ How to use this Makefile ------------------"
	@echo "make build   - Builds the executable."
	@echo "make clean   - Cleans the work directory."
	@echo "make help    - Show this help text."
	@echo "make test    - Runs the Unit tests and shows code coverage."
	@echo "make cover   - Generates HTML code coverage report."
	@echo "--------------------------------------------------------------"
