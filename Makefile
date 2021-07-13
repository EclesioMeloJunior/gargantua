.PHONY: gg
gg:
ifneq ($(wildcard ./bin),)
	@rm -rf "./bin"	
endif
	@go build -o ./bin/gg ./cmd/cli/...