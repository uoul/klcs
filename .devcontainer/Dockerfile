FROM node:22-alpine

# install golang
WORKDIR /
RUN GO_VERISON=1.25.6 \
    &&wget https://go.dev/dl/go$GO_VERISON.linux-amd64.tar.gz \
    && tar -xzf go$GO_VERISON.linux-amd64.tar.gz \
    && rm go$GO_VERISON.linux-amd64.tar.gz
ENV PATH=$PATH:/go/bin
RUN go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install golang.org/x/tools/gopls@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@latest

RUN npm install -g @angular/cli@21

# set workdir for project
WORKDIR /dev
