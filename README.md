# Microservices Makefile Guide

This document provides an overview of the commands included in the Makefile for managing your microservices. It covers how to run the services, generate protobuf files, build Docker images, and push them to a container registry.

## Prerequisites

Before using the Makefile, ensure you have the following installed:
- [Make](https://www.gnu.org/software/make/)
- [Docker](https://www.docker.com/)
- [Buf](https://buf.build/)

Ensure that your environment supports the `start` command for running processes in the background (e.g., Windows Command Prompt).

## Commands

### Run All Services

```bash
make run-all
```
This command starts all the microservices in the project:
- `auth-service`
- `order-service`
- `product-service`
- `gateway`

Each service is started in the background using the `start /b` command.

### Generate Protobuf Files

```bash
make buf
```
This command generates protobuf files using Buf. If any existing `*.go` files are found in the `pb` directory, they are deleted before generating the new files.

### Build Docker Images

```bash
make docker-build
```
This command builds Docker images for the following services:
- **API Gateway**: `vietquandeptrai/api-gateway`
- **Auth Service**: `vietquandeptrai/auth-svc`
- **Product Service**: `vietquandeptrai/product-svc`
- **Order Service**: `vietquandeptrai/order-svc`

### Push Docker Images to Registry

```bash
make docker-push
```
This command pushes the Docker images to the specified container registry:
- **API Gateway**: `vietquandeptrai/api-gateway`
- **Auth Service**: `vietquandeptrai/auth-svc`
- **Product Service**: `vietquandeptrai/product-svc`
- **Order Service**: `vietquandeptrai/order-svc`

## Notes

- Ensure you have logged in to the Docker registry before using `make docker-push`.
- The `start /b` command used for running services is specific to Windows environments. If you are on a Unix-based system, replace it with an equivalent command such as `&` for background execution.
- The Docker images use the `vietquandeptrai` namespace; update this if your registry uses a different naming convention.

## Example Usage

1. Generate protobuf files:
   ```bash
   make buf
   ```

2. Run all services:
   ```bash
   make run-all
   ```

3. Build Docker images:
   ```bash
   make docker-build
   ```

4. Push Docker images:
   ```bash
   make docker-push
   ```

