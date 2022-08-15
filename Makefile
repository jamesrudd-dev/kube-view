# Variables
RELEASE_TAG := ${git branch --show-current}

# Environmental Variables
.EXPORT_ALL_VARIABLES:
IN_PRODUCTION=false
KUBE_CONFIG_LOCATION=
WEB_SERVER_PATH=/kube-view
PUBLIC_URL=${WEB_SERVER_PATH}
IMAGE_TAG_FILTER=
NAMESPACE_FILTER=

run:
	make react-build && make go-run 

go-run:
	go run ./cmd/kube-view/*.go

go-install:
	go mod download

go-build:
	CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o build/kube-view cmd/kube-view/*.go

react-install:
	npm --prefix ./frontend install

react-build:
	npm --prefix ./frontend run build

react-development:
	npm --prefix ./frontend start

docker-build:
	docker build -t kube-view:$(RELEASE_TAG) .

docker-build-push:
	docker build -t jamesrudd/kube-view:$(RELEASE_TAG) .
	docker tag jamesrudd/kube-view:${RELEASE_TAG} jamesrudd/kube-view:latest
	docker push jamesrudd/kube-view:$(RELEASE_TAG)
	docker push jamesrudd/kube-view:latest
