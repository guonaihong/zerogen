all:
	# go build -o ./cmd/zero-gen ./cmd/zerogen.go
	go build ./cmd/zerogen/zerogen.go
	mv zerogen ../../test-zerogen
build-intel-mac:
	GOOS=darwin GOARCH=amd64 go build -o zerogen.intel ./cmd/zerogen/zerogen.go 
	mv zerogen.intel ../../test-zerogen
