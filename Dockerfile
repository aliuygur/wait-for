# # multi stage build, yo!
FROM golang:1.14
COPY . /app/
WORKDIR /app
RUN make deps build
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tcp-wait .

FROM scratch
WORKDIR /app
COPY --from=0 /app/bin/tcp-wait /app/tcp-wait
ENTRYPOINT [ "/app/tcp-wait" ]


# FROM golang:1.14
# COPY . /app/
# WORKDIR /app
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tcp-wait .

# FROM scratch
# COPY --from=0 /app/tcp-wait /app/tcp-wait
# ENTRYPOINT [ "/app/tcp-wait" ]
