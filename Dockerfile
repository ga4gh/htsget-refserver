FROM golang:latest

WORKDIR /usr/src/app

ENV SAMTOOLS_VERSION 1.9
ENV BCFTOOLS_VERSION 1.10.2

RUN apt-get update \
    && apt-get install --yes build-essential

RUN apt-get --yes install autoconf automake make gcc perl zlib1g-dev libbz2-dev liblzma-dev libcurl4-gnutls-dev libssl-dev libncurses5-dev

COPY go.mod go.sum ./
COPY cmd cmd
COPY internal internal
COPY data/gcp data/gcp
COPY data/config data/config
COPY deployments deployments
RUN mkdir temp
RUN go mod download

# Install samtools dependency
RUN cd /tmp \
    && wget https://github.com/samtools/samtools/releases/download/${SAMTOOLS_VERSION}/samtools-${SAMTOOLS_VERSION}.tar.bz2 \
    && tar xvjf samtools-${SAMTOOLS_VERSION}.tar.bz2 \
    && cd samtools-${SAMTOOLS_VERSION} \
    && ./configure --prefix=/usr/local \
    && make \
    && make install \
    && cd / && rm -rf /tmp/samtools-${SAMTOOLS_VERSION}

# Install bcftools dependency
RUN cd /tmp \
    && wget https://github.com/samtools/bcftools/releases/download/${BCFTOOLS_VERSION}/bcftools-${BCFTOOLS_VERSION}.tar.bz2 \
    && tar xvjf bcftools-${BCFTOOLS_VERSION}.tar.bz2 \
    && cd bcftools-${BCFTOOLS_VERSION} \
    && ./configure --prefix=/usr/local \
    && make \
    && make install \
    && cd / && rm -rf /tmp/bcftools-${BCFTOOLS_VERSION}

# Install htsget-refserver-utils dependency
RUN go mod download github.com/ga4gh/htsget-refserver-utils@v1.0.0 \
    && cd /go/pkg/mod/github.com/ga4gh/htsget-refserver-utils@v1.0.0 \
    && go install

ENV PATH="/usr/local:${PATH}"

RUN go build -o ./htsget-refserver ./cmd
EXPOSE 3000

CMD ["./htsget-refserver"]
