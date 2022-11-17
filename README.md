## Environment Setup

1. [Install required development software](.github/PREREQUISITES.md)
2. Start local minikube cluster `make start`
2. Add minikube ip to /etc/hosts (you will need administrator access) `make hosts`
3. Run `make skaffold`
4. Verification service will be available on http://verification-service.local
4. [Check make interact commands](.github/MAKE.md)

## Implementations

- [x] Development environment in Minikube. Helm and Skaffold used for deploy and reload
- [x] CQRS
- [x] Hexagonal architecture
- [x] DDD tactical design

## Use Cases

#### Verification
- [x] Create
- [x] Approve
- [x] Decline
- [x] Get by uuid

## Stack

- Golang 1.19
- PostgreSQL 15
