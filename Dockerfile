FROM golang as builder

RUN mkdir /build
WORKDIR /build 

RUN export GO111MODULE=on

RUN go install github.com/bunkieproject/bunkie_be@latest
RUN cd /build && git clone https://github.com/bunkieproject/bunkie_be.git

RUN cd /build/bunkie_be && go build -o bunkie_be

EXPOSE 8080

CMD ["/build/bunkie_be/bunkie_be"]