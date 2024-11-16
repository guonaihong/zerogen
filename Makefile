all:
	# go build -o ./cmd/zero-gen ./cmd/zerogen.go
	go build ./cmd/zerogen.go
	mv zerogen ../../test-zerogen
