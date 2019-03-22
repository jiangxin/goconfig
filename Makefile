TARGET=goconfig

define message
	@echo "### $(1)"
endef

all: $(TARGET)

goconfig: $(shell find . -name '*.go')
	go build -o $@ cmd/goconfig/main.go

test:
	$(call message,Testing goconfig using golint for coding style)
	@golint
	$(call message,Testing goconfig for unit tests)
	@go test

clean:
	rm -f $(TARGET)
