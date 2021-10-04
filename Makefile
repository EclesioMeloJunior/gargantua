.PHONY: gg
gg:
ifneq ($(wildcard ./bin),)
	@rm -rf "./bin"	
endif
	@go1.17.1 build -o ./bin/gg ./cmd/cli/...