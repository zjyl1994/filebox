build:
	CGO_ENABLE=0 go build -o filebox .
mini:
	CGO_ENABLE=0 go build --ldflags="-s -w" -o filebox . && upx filebox
