FROM --platform=${BUILDPLATFORM} node:alpine AS build-stage

# get target platform
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

# get target platform
ARG TARGETOS
ARG TARGETARCH

# install golang
WORKDIR /
RUN GO_VERSION=1.25.6 \
    && wget https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz \
    && tar -xzf go$GO_VERSION.linux-amd64.tar.gz \
    && rm go$GO_VERSION.linux-amd64.tar.gz
ENV PATH=$PATH:/go/bin

RUN npm install -g @angular/cli

# set workdir for project
WORKDIR /app
COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} cd ./backend && go mod download && go build -o klcs main.go
RUN cd frontend/klcs && npm i &&  ng build --configuration production

# Deploy the application binary into a lean image
FROM alpine:latest AS build-release-stage

WORKDIR /app

COPY --from=build-stage /app/backend/klcs /app
COPY --from=build-stage /app/frontend/klcs/dist/klcs/browser /app/wwwroot

ENTRYPOINT ["/app/klcs"]
