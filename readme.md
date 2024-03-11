### Sobre o projeto
Esse é uma pequena aplicação que permite realizar teste de stress.

### Pré Requisitos
Go 1.22.1

### Como executar?
```
go run main.go --url=http://google.com --requests=1000 --concurrency=10
```

### Executando via Docker?
```
docker run thedevices/go-stress-test:latest --url=http://www.google.com --requests=100 --concurrency=10 
```