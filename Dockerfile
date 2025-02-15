FROM node:22-alpine AS build-stage

# install golang
WORKDIR /
RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz \
    && tar -xzf go1.24.0.linux-amd64.tar.gz \
    && rm go1.24.0.linux-amd64.tar.gz
ENV PATH=$PATH:/go/bin

RUN npm install -g @angular/cli

# set workdir for project
WORKDIR /app
COPY . .
RUN cd backend/core && go mod download && go build -o klcs main.go
RUN cd frontend/klcs && ng build --configuration production

# Deploy the application binary into a lean image
FROM alpine:latest AS build-release-stage

WORKDIR /app

COPY --from=build-stage /app/backend/core/klcs /app
COPY --from=build-stage /app/frontend/klcs/dist/klcs/browser /app/wwwroot

ENTRYPOINT ["/app/klcs"]
