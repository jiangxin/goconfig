TARGET=goconfig

define message
	@echo "### $(1)"
endef

test: $(TARGET) golint
	$(call message,Testing goconfig using golint for coding style)
	@golint
	$(call message,Testing goconfig for unit tests)
	@go test

golint:
	@if ! type golint >/dev/null 2>&1; then \
		go get golang.org/x/lint/golint; \
	fi

all: $(TARGET)

goconfig: $(shell find . -name '*.go')
	$(call message,Build $@)
	@go build -o $@ cmd/goconfig/main.go

clean:
	rm -f $(TARGET)

.PHONY: test clean golint
