build:
	go build -o filebox .
mini:
	go build --ldflags="-s -w" -o filebox . && upx filebox
