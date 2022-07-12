run:
	npm --prefix ./frontend/my-app run build && go run main.go

react-build:
	npm --prefix ./frontend/my-app run build

react-run-go:
	npm --prefix ./frontend/my-app run server

react-start:
	npm --prefix ./frontend/my-app start
