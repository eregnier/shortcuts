dev:
	export DEV=1 && air

build:
	go build main.go

install: build
	chmod +x main.go
	sudo mv main /usr/local/bin/my-actions
