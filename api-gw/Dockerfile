
FROM --platform=linux/amd64 golang:1.19 as builder

ARG API_GW_COMPONENT_NAME 
WORKDIR /app
COPY . .
WORKDIR /app/api-gw
RUN go mod tidy && go mod download && GOOS=linux GOARCH=amd64 go build -buildvcs=false -o bin/${API_GW_COMPONENT_NAME} 

FROM gcr.io/nsx-sm/photon:4.0
WORKDIR /bin
COPY --from=builder /app/api-gw/bin/api-gw .
USER 65532:65532
ENTRYPOINT "/bin/api-gw"
