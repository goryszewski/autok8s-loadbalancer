#VAR

REPO=repo.internal:5000
IMAGE=loadbalancer-external
TAG=latest

#TARGET

build:
	DOCKER_BUILDKIT=0 docker build -t $(REPO)/$(IMAGE):$(TAG) .
run:
	go run ./cmd/main.go
publish:
	docker push $(REPO)/$(IMAGE):$(TAG)