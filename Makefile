all: go

exec = dota-predictor

go: main.go
	go build -o bin/$(exec) main.go

clean:
	rm -f bin/$(exec) *~ *#
	rm -rf pkg
	rm -rf src/gopkg.in

deps:
	go get github.com/gorilla/mux
	go get -u gorm.io/gorm
	go get -u gorm.io/driver/sqlite


re: clean all