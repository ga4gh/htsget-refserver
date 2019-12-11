FROM golang:latest

WORKDIR /usr/src/app

ENV SAMTOOLS_VERSION 1.10

RUN apt-get update \
    && apt-get install --yes build-essential

RUN apt-get --yes install autoconf automake make gcc perl zlib1g-dev libbz2-dev liblzma-dev libcurl4-gnutls-dev libssl-dev libncurses5-dev

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN cd /tmp \
    && wget https://github.com/samtools/samtools/releases/download/${SAMTOOLS_VERSION}/samtools-${SAMTOOLS_VERSION}.tar.bz2 \
    && tar xvjf samtools-${SAMTOOLS_VERSION}.tar.bz2 \
    && cd samtools-${SAMTOOLS_VERSION} \
    && ./configure --prefix=/usr/local \
    && make \
    && make install \
    && cd / && rm -rf /tmp/samtools-${SAMTOOLS_VERSION}

ENV PATH="/usr/local:${PATH}"

RUN go build -o htsref ./cmd
EXPOSE 3000

CMD ["./htsref"]
