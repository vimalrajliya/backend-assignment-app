from golang:1.24.3
WORKDIR /app

RUN go install github.com/air-verse/air@latest
ENV PATH="/go/bin:${PATH}"
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]
