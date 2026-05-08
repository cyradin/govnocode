ARG GO_VERSION=1.26
FROM golang:${GO_VERSION}-trixie

ARG GOPRIVATE=
ENV GOPRIVATE=${GOPRIVATE}

RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    bash-completion \
    wget \
    git \
    nano \
    vim \
    curl \
    net-tools \
    coreutils \
    sudo \
    build-essential \
    jq \
    procps \
    ripgrep \
    fd-find \
    tree \
    less \
    tini

RUN apt-get -y autoremove \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN ln -s /usr/bin/fdfind /usr/local/bin/fd

RUN addgroup user \
    && adduser --disabled-password --shell=/bin/bash --ingroup=user user \
    && echo 'user ALL=NOPASSWD: ALL' > /etc/sudoers.d/user \
    && chmod 0440 /etc/sudoers.d/user \
    && visudo --check

RUN mkdir -p /app && chown user:user /app

USER user

ENV HOME=/home/user

WORKDIR /app

RUN git config --global --add safe.directory '*'

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["bash", "-c", "trap : TERM INT; sleep infinity & wait"]