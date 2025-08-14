#--------------------------------------------------------------
# base build image
#--------------------------------------------------------------
FROM mcr.microsoft.com/vscode/devcontainers/go:1.24
#--------------------------------------------------------------
# ARG
#--------------------------------------------------------------
# default
ARG USER=vscode
ARG WORKDIR=/workspace
# version(base)
ARG NODE_VERSION=24
ARG GITLEAKS_VERSION=8.18.0
ARG GH_VERSION=2.73.0
# version(go)
ARG GOLANGCI_LINT_VERSION=v2.3.1

# x86_64
# ARG ARCH_BASE=amd64
# ARG ARCH_BASE2=x86_64
# ARG ARCH_GITLEAKS=x64
# ARG ARCH_AWS_CLI=x86_64
# ARG ARCH_SESSION_MANAGER_PLUGIN=64bit
# ARG ARCH_KO=x86_64

# arm64
ARG ARCH_BASE=arm64
ARG ARCH_BASE2=aarch_64
ARG ARCH_GITLEAKS=arm64
ARG ARCH_AWS_CLI=aarch64
ARG ARCH_SESSION_MANAGER_PLUGIN=arm64
ARG ARCH_KO=arm64

#--------------------------------------------------------------
# ENV
#--------------------------------------------------------------
ENV CGO_ENABLED=0
ENV GOOS=linux
# x86_64
# ENV GOARCH=amd64
# arm64
ENV GOARCH=arm64
ENV GOPATH=/go

ENV PATH="/usr/local/bin/gh_${GH_VERSION}_linux_${ARCH_BASE}/bin:${PATH}"

#--------------------------------------------------------------
# workdir
#--------------------------------------------------------------
RUN mkdir -p ${WORKDIR}
WORKDIR  ${WORKDIR}

#--------------------------------------------------------------
# Install software-properties-common
#--------------------------------------------------------------
RUN apt update && \
    apt -y install software-properties-common

#--------------------------------------------------------------
# Install dependent packages
#--------------------------------------------------------------
RUN apt update && \
    apt -y install make openssl git tar curl zip jq groff gcc bison python3 python3-dev python3-pip fish shellcheck shfmt

# Install Node.js
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash -E - && \
    apt-get install -y nodejs=${NODE_VERSION}.0.0-1nodesource1

# Install pre-commit
RUN pip3 install pre-commit --break-system-packages --ignore-installed distlib

# Install gitleaks
RUN curl -sLO https://github.com/gitleaks/gitleaks/releases/download/v${GITLEAKS_VERSION}/gitleaks_${GITLEAKS_VERSION}_linux_${ARCH_GITLEAKS}.tar.gz && \
    mkdir -p ./gitleaks_${GITLEAKS_VERSION}_linux_${ARCH_GITLEAKS} && \
    tar -zxvf gitleaks_${GITLEAKS_VERSION}_linux_${ARCH_GITLEAKS}.tar.gz -C ./gitleaks_${GITLEAKS_VERSION}_linux_${ARCH_GITLEAKS} && \
    chmod +x ./gitleaks_${GITLEAKS_VERSION}_linux_${ARCH_GITLEAKS}/gitleaks && \
    mv ./gitleaks_${GITLEAKS_VERSION}_linux_${ARCH_GITLEAKS}/gitleaks /usr/local/bin/gitleaks && \
    rm -rf gitleaks_${GITLEAKS_VERSION}_linux_${ARCH_GITLEAKS}.tar.gz ./gitleaks_${GITLEAKS_VERSION}_linux_${ARCH_GITLEAKS}

# Install gh
RUN wget https://github.com/cli/cli/releases/download/v${GH_VERSION}/gh_${GH_VERSION}_linux_${ARCH_BASE}.tar.gz && \
    tar -xvf gh_${GH_VERSION}_linux_${ARCH_BASE}.tar.gz && \
    cp -riv gh_${GH_VERSION}_linux_${ARCH_BASE} /usr/local/bin/ && \
    rm -rf gh_${GH_VERSION}_linux_${ARCH_BASE}.tar.gz gh_${GH_VERSION}_linux_${ARCH_BASE}

# Install npm packages
RUN npm install -g markdownlint-cli markdown-link-check cspell

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
# Install AWS CLI
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-${ARCH_AWS_CLI}.zip" -o "awscliv2.zip" && \
    unzip awscliv2.zip && \
    ./aws/install && \
    rm -rf awscliv2.zip

# Install aws-mfa
RUN  pip3 install --break-system-packages aws-mfa && \
    # Install awsp
    npm install -g awsp && \
    # echo 'alias awsp="source _awsp"' >> /home/${USER}/.bashrc && \
    # Install aws-sam-cli
    # curl -sSOL https://github.com/aws/aws-sam-cli/releases/latest/download/aws-sam-cli-linux-x86_64.zip && \
    # unzip aws-sam-cli-linux-x86_64.zip -d sam-installation && \
    # ./sam-installation/install && \
    # Install Session Manager Plugin
    curl "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_${ARCH_SESSION_MANAGER_PLUGIN}/session-manager-plugin.deb" -o "session-manager-plugin.deb" && \
    dpkg -i session-manager-plugin.deb

#--------------------------------------------------------------
# Install go related
#--------------------------------------------------------------
RUN go install -v golang.org/x/tools/gopls@latest && \
    go install -v golang.org/x/tools/...@latest && \
    go install -v github.com/cweill/gotests/...@latest && \
    go install -v honnef.co/go/tools/...@latest && \
    go install -v golang.org/x/lint/golint@latest && \
    go install -v github.com/haya14busa/goplay/cmd/goplay@latest && \
    go install -v github.com/go-delve/delve/cmd/dlv@latest && \
    go install -v github.com/sqs/goreturns@latest && \
    go install -v github.com/josharian/impl@latest && \
    go install -v github.com/davidrjenni/reftools/cmd/fillstruct@latest && \
    go install -v github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest && \
    go install -v github.com/ramya-rao-a/go-outline@latest && \
    go install -v github.com/acroca/go-symbols@latest && \
    go install -v github.com/godoctor/godoctor@latest && \
    go install -v github.com/rogpeppe/godef@latest && \
    go install -v github.com/zmb3/gogetdoc@latest && \
    go install -v github.com/fatih/gomodifytags@latest && \
    go install -v github.com/mgechev/revive@latest && \
    # check go vulnerabilities
    go install -v golang.org/x/vuln/cmd/govulncheck@latest && \
    # hot reload tool: air
    # go install -v github.com/cosmtrek/air@latest && \
    # hot reload tool: reflex
    go install -v github.com/cespare/reflex@latest && \
    # generate mock
    go install -v github.com/golang/mock/mockgen@latest && \
    # go swagger tool
    go install -v github.com/go-swagger/go-swagger/cmd/swagger@latest && \
    go install -v github.com/swaggo/swag/cmd/swag@latest && \
    # Coverage outputs covertura format.
    go install -v github.com/t-yuki/gocover-cobertura@latest 2>&1 && \
    # Install golangci-lint
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}

RUN chown -R vscode:golang /go/pkg

#--------------------------------------------------------------
# Clean up
#--------------------------------------------------------------
RUN rm -rf /tmp/tmp.* && \
    apt autoremove -y && \
    apt clean -y && \
    rm -rf /var/lib/apt/lists/*

ENTRYPOINT []
CMD []
