name = setlx2python_playground

package = .

all: build run

run:
	./$(name) -port="8080" -mode="dev"

build:
	GO111MODULE=on go build -o $(name) $(package)

release:
	CGO_ENABLED=0 GO111MODULE=on go build -ldflags="-s -w" -a -installsuffix nocgo -o setlxplay $(package)