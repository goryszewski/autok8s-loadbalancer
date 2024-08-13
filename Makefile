#VAR

REPO=repo.internal:5000
IMAGE=loadbalancer-external
TAG=latest

#TARGET

update:
	go get -u github.com/goryszewski/libvirtApi-client

build:
	DOCKER_BUILDKIT=0 docker build -t $(REPO)/$(IMAGE):$(TAG) .

run:
	go run ./cmd/main.go

dev:
	go run ./cmd/main.go -dev

publish: build
	docker push $(REPO)/$(IMAGE):$(TAG)