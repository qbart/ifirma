build:
	mkdir -p bin/
	go build -o bin/ifirma

install: build
	cp bin/ifirma /usr/local/bin
