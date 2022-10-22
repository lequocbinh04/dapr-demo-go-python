FROM golang:1.18 as builder
ENV GOOS linux
ENV CGO_ENABLED 0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app

FROM golang:1.18 as production
#RUN apk add --no-cache ca-certificates
COPY --from=builder app .
#RUN apk update && apk add wget
# Install  daprd
RUN wget -q https://raw.githubusercontent.com/dapr/cli/master/install/install.sh -O - | /bin/bash
ARG DAPR_BUILD_DIR
COPY $DAPR_BUILD_DIR /opt/dapr
ENV PATH="/opt/dapr/:${PATH}"
RUN dapr init --slim

EXPOSE 8080
EXPOSE 3030
ENTRYPOINT ["dapr"]
CMD ["run", "--app-id", "server", "--app-port", "4999", "--dapr-grpc-port", "3030", "./app"]
