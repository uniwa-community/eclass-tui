GOFLAGS = -v
RUNFLAGS =

PREFIX = /usr/local/

all: options build

options:
	@echo "compile options:"
	@echo "GOFLAGS  = $(GOFLAGS)"

build: tidy
	cd src && \
	go build $(GOFLAGS) -o ../eclass-tui

tidy:
	cd src/ && \
	go mod tidy
	touch tidy

test:
	cd src/ && \
	go test .

clean:
	rm -f ./eclass-tui ./tidy

run: build
	./eclass-tui $(RUNFLAGS)

install: build
	@echo "Installing into $(DESTDIR)$(PREFIX)bin/.."
	cp -f ./eclass-tui $(DESTDIR)$(PREFIX)bin/eclass-tui
	chmod 755          $(DESTDIR)$(PREFIX)bin/eclass-tui

uninstall:
	@echo "Uninstalling from $(DESTDIR)$(PREFIX)bin/.."
	rm -f $(DESTDIR)$(PREFIX)bin/eclass-tui

.PHONY: run build clean install uninstall options all test # not tidy
