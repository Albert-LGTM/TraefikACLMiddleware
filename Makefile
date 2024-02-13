.PHONY: build run clean

build:
	go build -o myapp

run:
	./myapp

clean:
	rm -f myapp

