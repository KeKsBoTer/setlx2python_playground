FROM golang:1.12 as builder
WORKDIR /server/

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY  cmd cmd
COPY  execute.go .
COPY  fetch.go .
COPY  handler.go .
COPY  router.go .
COPY  transpile.go .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -a -installsuffix nocgo -o playground github.com/keksboter/setlx2python_playground/cmd/setlx2python_playground

FROM python:3.6
WORKDIR /root/
COPY www www
COPY --from=builder /server/playground .
RUN pip install git+https://github.com/KeKsBoTer/setlx2python.git
ENTRYPOINT [ "./playground","-mode","prod","-port","80"]
EXPOSE 80