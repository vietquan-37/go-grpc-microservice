PHONY: run-all  buf
run-all:

	cd auth-service && start /b make run-server
	cd order-service && start /b make run-server
	cd product-service && start /b  make run-server
buf:
	@if exist pb\*.go del /Q pb\*.go
	buf generate

