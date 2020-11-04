# Config
VERSION=1.0.0
NAME=bcache
BUILD_TIME=`date +%FT%T%z`
BUILD_ID=`git rev-parse --short=10 HEAD`

###################################
# DO NOT MODIFY BELLOW THIS POINT #
###################################

BUILD=go build -ldflags="\
	-s \
	-w \
	-X main.Version=${VERSION} \
	-X main.BuildId=${BUILD_ID} \
	-X main.BuildDate=${BUILD_TIME} \
"

default: linux64

linux64:
	GOOS=linux GOARCH=amd64 $(BUILD) -o ${NAME}

darwin64:
	GOOS=darwin GOARCH=amd64 $(BUILD) -o ${NAME}

windows64:
	GOOS=windows GOARCH=amd64 $(BUILD) -o ${NAME}

all:
	GOOS=linux GOARCH=amd64 $(BUILD) -o ${NAME}_linux64
	GOOS=darwin GOARCH=amd64 $(BUILD) -o ${NAME}_darwin64
	GOOS=windows GOARCH=amd64 $(BUILD) -o ${NAME}_windows64

clean:
	@echo "Nothing to do"
