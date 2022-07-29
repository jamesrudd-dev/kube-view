run:
	npm --prefix ./frontend run build && go run main.go

react-build:
	npm --prefix ./frontend run build

react-run-go:
	npm --prefix ./frontend run server

react-start:
	npm --prefix ./frontend start
