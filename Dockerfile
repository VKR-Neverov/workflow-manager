FROM golang:alpine

COPY ./configs ./configs
COPY bin/main /main

CMD ["/main"]