NOVENDOR_PATH = $$(glide novendor)
.PHONY: test

glide:
	-rm glide.lock
	-rm -r vendor
	glide cache-clear
	glide install

test:
	go clean

run:
	go run main.go
