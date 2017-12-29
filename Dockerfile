# multi stage build, yo!
FROM golang:1.8
COPY . /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wait .

FROM scratch
COPY --from=0 /app/wait /app/wait
ENTRYPOINT [ "/app/wait" ]