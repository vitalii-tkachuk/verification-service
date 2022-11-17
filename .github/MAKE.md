### Start minikube cluster (can be used for further runs after init)
```bash
make start
```

### Stop minikube cluster
```bash
make stop
```

### Debug local helm config
```bash
make helm-debug
```

### Upgrade local minikube cluster
```bash
make helm-upgrade
```

### Shell to postgres pod
```bash
make postgres
```

### List logs for application pod
```bash
make logs
```

### Run skaffold in development environment
```bash
make skaffold
```

### Add minikube ip to /etc/hosts file (need to be executed only once during initial setup)
```bash
make hosts
```

### Format application code using `go fmt`
```bash
make format
```

### Vet application code using `go vet`
```bash
make vet
```

### Lint application code using `golangci-lint`
```bash
make lint
```

### Run application tests
```bash
make test
```
