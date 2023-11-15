VERSION = 0.3

build:
	GOPATH=/home/vagrant/go GOROOT=/usr/local/go go build -o app main.go

image:build
	docker build . --tag localhost:5001/http:$(VERSION)
	docker push localhost:5001/http:$(VERSION)

clean:
	-kubectl delete -f deploy.yaml
	-kubectl delete cm simpleconf

deploy:clean
	kubectl create configmap simpleconf --from-file=data=./conf/xx.yaml
	kubectl apply -f deploy.yaml
	kubectl get pods
