
NOW=`date -u '+%Y-%m-%d %H:%M:%S'`
BUILD_NUMBER_FILE=build.txt
PORT_NUMBER=8014
SOURCES=$(shell find . -name '*.go')

build: $(SOURCES)
	@if ! test -f $(BUILD_NUMBER_FILE); then echo 0 > $(BUILD_NUMBER_FILE); fi
	@echo $$(($$(cat $(BUILD_NUMBER_FILE)) + 1)) > $(BUILD_NUMBER_FILE)

	@go build -ldflags "-X 'github.com/antsy/tournament/utils.Buildtime=${NOW}' -X 'github.com/antsy/tournament/utils.Buildnumber=`cat $(BUILD_NUMBER_FILE)`'"
	@echo "built the tournament server for the `cat $(BUILD_NUMBER_FILE)`th time"

kill:
	fuser -k ${PORT_NUMBER}/tcp

run:
	./tournament &

test:
	@go test ./tests -parallel 8

.PHONY: build kill run test
