TARGET=goconfig

all: $(TARGET)

goconfig: $(shell find . -name '*.go')
	go build -o $@ cmd/goconfig/main.go

test:
	golint && go test ./...

clean:
	rm -f $(TARGET)
