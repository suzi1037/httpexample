VERSION = 1.2.3

build:
	GOPATH=/home/vagrant/go GOROOT=/usr/local/go go build -o app main.go

image:build
	docker build . --tag localhost:5001/httpserver:$(VERSION)
	docker push localhost:5001/httpserver:$(VERSION)

deploy:
	helm install foo ./helm
	kubectl get pods

helm-dry:
	helm install foo ./helm --dry-run

