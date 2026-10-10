.PHONY: build test
COMMON_FLAGS := --username admin --password Harbor12345 --baseurl http://localhost:30002
PROGR := ./go-harbor-cli

build:
	go build -o go-harbor-cli

create-project: build
	$(PROGR) project create $(COMMON_FLAGS) --name deleteme --quota-size-gib 4 --admin-users "pippo,pluto"

get-quota: build
	$(PROGR) quota get $(COMMON_FLAGS) --project-id 36

set-quota: build
	$(PROGR) quota set $(COMMON_FLAGS) --project-id 36 --quota-size-gib 7