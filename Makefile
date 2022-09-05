OBJ = md2vim.go vimdoc.go
BIN = md2vim
BDIR = /usr/local/bin

build: $(OBJ)
	go build

.PHONY: install
install: $(BIN)
	mv $(BIN) $(BDIR)
