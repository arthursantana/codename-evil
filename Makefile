all: dist/evil dist/static

clean:
	rm -rf dist/evil dist/static

dist/evil: server/*
	cd server; go build -o ../dist/evil; cd ..

dist/static: node_modules gulpfile.js client/* client/*/*
	node node_modules/gulp/bin/gulp.js

node_modules: package.json
	npm install

.PHONY: all clean init

init:
	go get github.com/gorilla/websocket

run:
	cd dist; ./evil ${ARGS}; cd ..

watch:
	node node_modules/gulp/bin/gulp.js watch

