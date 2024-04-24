FROM golang:1.21.4

WORKDIR /registration

#workaround, need this because some scripts are generate via go generate
# RUN git clone https://github.com/PelicanPlatform/pelican.git &&\
#     cd pelican && \
#     git checkout tags/v7.2.0 && \
#     make generate

COPY go.mod go.mod
COPY go.sum go.sum
COPY controllers controllers
COPY models models
COPY utilities utilities
COPY input.txt input.txt

RUN go mod download
COPY main.go main.go
RUN go get registration

RUN go build -o main .

ENTRYPOINT [ "./main" ]
# ENTRYPOINT ["sleep", "infinity"]