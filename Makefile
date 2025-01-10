PHONY: run-all  buf
run-all:
	cd auth-service && start /b make run-server
	cd order-service && start /b make run-server
	cd product-service && start /b  make run-server
	cd gateway && start /b make run-server
buf:
	@if exist pb\*.go del /Q pb\*.go
	buf generate
docker-build:
	cd gateway && docker build -t vietquandeptrai/api-gateway .
	docker build -f ./auth-service/Dockerfile -t  vietquandeptrai/auth-svc .
	docker build -f ./product-service/Dockerfile -t vietquandeptrai/product-svc .
	docker build -f ./order-service/Dockerfile -t vietquandeptrai/order-svc .
docker-push:
	cd gateway && docker push vietquandeptrai/api-gateway
	cd auth-service && docker push vietquandeptrai/auth-svc
	cd product-service && docker push vietquandeptrai/product-svc
	cd order-service && docker push vietquandeptrai/order-svc



