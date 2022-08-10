RELEASE_TAG := release-1.0

go-run:
	go run ./cmd/kube-view/*.go

react-build:
	npm --prefix ./frontend run build

react-start:
	npm --prefix ./frontend start

docker-build:
	docker build -t jamesrudd/kube-view:$(RELEASE_TAG) .
	docker tag jamesrudd/kube-view:${RELEASE_TAG} jamesrudd/kube-view:latest

docker-push:
	docker push jamesrudd/kube-view:$(RELEASE_TAG)
	docker push jamesrudd/kube-view:latest

docker-build-push:
	docker build -t jamesrudd/kube-view:$(RELEASE_TAG) .
	docker tag jamesrudd/kube-view:${RELEASE_TAG} jamesrudd/kube-view:latest
	docker push jamesrudd/kube-view:$(RELEASE_TAG)
	docker push jamesrudd/kube-view:latest
