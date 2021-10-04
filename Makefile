.PHONY: gg
gg:
ifneq ($(wildcard ./bin),)
ifeq ($(OS),Windows_NT)
	@rmdir /s /q "./bin"
else
	@rm -rf "./bin"	
endif
endif
	@go build -o ./bin/gg ./cmd/cli/...
