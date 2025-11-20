#--------------------------------------------------------------
# base build image
#--------------------------------------------------------------
FROM mcr.microsoft.com/vscode/devcontainers/go:1.25
#--------------------------------------------------------------
# ARG
#--------------------------------------------------------------
# default
ARG USER=vscode
ARG WORKDIR=/workspace
# version(base)
ARG NODE_VERSION=24 # https://github.com/nodejs/node/releases
ARG AQUA_VERSION=v2.55.1 # https://github.com/aquaproj/aqua/releases

# x86_64
# ARG ARCH_BASE=amd64
# ARG ARCH_BASE2=x86_64
# ARG ARCH_SESSION_MANAGER_PLUGIN=64bit
# ARG ARCH_KO=x86_64

# arm64
ARG ARCH_BASE=arm64
ARG ARCH_BASE2=aarch_64
ARG ARCH_SESSION_MANAGER_PLUGIN=arm64
ARG ARCH_KO=arm64

#--------------------------------------------------------------
# ENV
#--------------------------------------------------------------
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOPATH=/go
# x86_64
# ENV GOARCH=amd64
# arm64
ENV GOARCH=arm64

#--------------------------------------------------------------
# workdir
#--------------------------------------------------------------
RUN mkdir -p ${WORKDIR}
WORKDIR  ${WORKDIR}

#--------------------------------------------------------------
# Install minimal apt/https tooling (avoid software-properties-common)
#--------------------------------------------------------------
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    ca-certificates \
    gnupg \
    dirmngr \
    apt-transport-https \
    lsb-release

#--------------------------------------------------------------
# Install dependent packages
#--------------------------------------------------------------
RUN apt update && \
    apt -y install make openssl git tar curl zip jq groff gcc bison python3 python3-dev python3-pip fish shellcheck shfmt bats

# Install Node.js
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash -E - && \
    apt-get install -y nodejs=${NODE_VERSION}.0.0-1nodesource1

# Install npm packages
RUN npm install -g markdownlint-cli markdown-link-check cspell textlint textlint-rule-preset-ja-technical-writing textlint-rule-no-dead-link

# Install aqua
RUN go install github.com/aquaproj/aqua/v2/cmd/aqua@${AQUA_VERSION}

#--------------------------------------------------------------
# Install for MCP Server
#--------------------------------------------------------------
# Install for vscode user
USER ${USER}
RUN curl -LsSf https://astral.sh/uv/install.sh | sh

# Switch back to root for remaining operations
USER root

#--------------------------------------------------------------
# AWS environment
#--------------------------------------------------------------
# Install aws-mfa
RUN pip3 install --break-system-packages aws-mfa && \
    # Install awsp
    npm install -g awsp && \
    # Install Session Manager Plugin
    curl "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_${ARCH_SESSION_MANAGER_PLUGIN}/session-manager-plugin.deb" -o "session-manager-plugin.deb" && \
    dpkg -i session-manager-plugin.deb

# Install awsdac
RUN apt update && \
    apt -y install fonts-liberation && \
    go install github.com/awslabs/diagram-as-code/cmd/awsdac@latest && \
    go install github.com/awslabs/diagram-as-code/cmd/awsdac-mcp-server@latest

#--------------------------------------------------------------
# Install go related
#--------------------------------------------------------------
    # generate mock
RUN go install -v github.com/golang/mock/mockgen@latest && \
    # Coverage outputs covertura format.
    go install -v github.com/t-yuki/gocover-cobertura@latest 2>&1

RUN chown -R vscode:golang /go/pkg && \
    chmod -R g+w /go

#--------------------------------------------------------------
# Clean up
#--------------------------------------------------------------
RUN rm -rf /tmp/tmp.* && \
    apt autoremove -y && \
    apt clean -y && \
    rm -rf /var/lib/apt/lists/*

ENTRYPOINT []
CMD []
