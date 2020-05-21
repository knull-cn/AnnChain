# compile environment;
FROM annchain/builder:go1.12 as builder
RUN apk add leveldb-dev
#copy files;
ADD . /AnnChain
WORKDIR /AnnChain
RUN GO111MODULE="on" GOPROXY="https://goproxy.io" make genesis

# package environment;
FROM annchain/runner:alpine3.11
RUN apk add leveldb
WORKDIR /genesis
COPY --from=builder /AnnChain/build/genesis /bin/
#COPY build/genesis /bin/
ENTRYPOINT ["genesis"]

