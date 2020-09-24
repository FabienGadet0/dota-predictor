all: go

exec = dota-predictor

go: src/*.go
	go build -o bin/$(exec) src/*.go

clean:
	rm -f bin/$(exec) *~ *#
	rm -rf pkg
	rm -rf src/gopkg.in

deps:
	go get github.com/gorilla/mux

re: clean all