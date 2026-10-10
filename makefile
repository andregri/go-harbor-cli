.PHONY: build test
build:
	go build -o go-harbor-cli

test: build
	./go-harbor-cli project create --name deleteme --username admin --password Harbor12345 --baseurl http://localhost:30002 --quota-size-gib 4 --admin-members "pippo,pluto"