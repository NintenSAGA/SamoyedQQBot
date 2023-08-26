GO = go
SRC = ./src
OUT = ./out
BIN = $(OUT)/bot

.PHONY: build clean

build:
	mkdir -p $(OUT)
	$(GO) build -o $(BIN) $(SRC)

clean:
	rm $(BIN)

run: $(BIN) ./data/wordlist.txt
	WORDLIST_PATH="./data/wordlist.txt" $(BIN)
