# Usar a imagem oficial do Go como imagem de construção
FROM golang:1.22.1-alpine as builder

# Definir o diretório de trabalho dentro do container
WORKDIR /app

# Copiar o módulo Go para o diretório de trabalho
COPY . .

# Baixar dependências. Dependências serão armazenadas no cache se o go.mod e o go.sum não forem alterados
RUN go mod download

# Compilar o aplicativo Go para um binário estático
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Usar a imagem scratch como imagem final para um tamanho menor da imagem
FROM scratch

WORKDIR /root/

# Copiar o binário estático para a imagem final
COPY --from=builder /app/main .

# Executar o binário
ENTRYPOINT ["/root/main"]
