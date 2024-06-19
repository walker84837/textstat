operatingsystem := $(shell go env GOOS)

.PHONY: all clean release debug

all:
	@echo "Building release for '$(operatingsystem)'..."
	go build -ldflags="-s -w -linkmode external" -o tstat src/textstat.go 
clean:
	@rm -rf ./*.a *.o core *.log tstat
