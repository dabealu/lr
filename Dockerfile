FROM       golang:1.8-alpine
COPY       . /usr/src/lr
RUN        go build -o /usr/local/bin/lr /usr/src/lr/lr.go
ENTRYPOINT ["/usr/local/bin/lr"]
CMD        ["help"]
