.PHONY: gg
gg:
ifneq ($(wildcard ./bin),)
	@rd /s /q "./bin"	
endif
	@go build -o ./bin/gg.exe ./cmd/cli/...