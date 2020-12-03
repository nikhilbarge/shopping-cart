# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
 
BINARY_NAME=shopping-cart
DOCKER_IMAGE=shopping-cart

# For Development Build #################################################################
# Docker.io username and tag
DOCKER_USER=nikhilbarge
DOCKER_IMAGE_TAG=test
# For Development Build #################################################################

 

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test: 
	$(GOTEST) -v ./...
  
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

modverify:
	$(GOMOD) verify

modtidy:
	$(GOMOD) tidy

moddownload:
	$(GOMOD) download
  
docker-build:
	docker build -t $(DOCKER_USER)/$(DOCKER_IMAGE):$(DOCKER_IMAGE_TAG) -f Dockerfile .
 

docker-push:
	docker push $(DOCKER_USER)/$(DOCKER_IMAGE):$(DOCKER_IMAGE_TAG)
	 

buildlocal: build docker-build clean

all: build docker-build docker-push 