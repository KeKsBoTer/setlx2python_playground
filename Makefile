name = setlx2python_playground

package = .

all: build run

run:
	./$(name).exe -port="8080" -mode="dev"

build:
	go build -o $(name).exe $(package)/cmd/$(name)

release:
	CGO_ENABLED=0 GO111MODULE=on go build -ldflags="-s -w" -a -installsuffix nocgo -o setlxplay $(package)