.PHONY: build test install uninstall clean

BINARY := todome
INSTALL_DIR := $(HOME)/.bin

build:
	go build -o $(BINARY) .

test:
	go test ./...

install: build
	@mkdir -p $(INSTALL_DIR)
	cp $(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "$(BINARY) を $(INSTALL_DIR) にインストールした 🗡️"

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY)
	@echo "$(BINARY) をアンインストールした"

clean:
	rm -f $(BINARY)
