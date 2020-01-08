FROM scratch
WORKDIR /go/src/myapp
COPY go.mod go.sum ./
ADD main /
COPY . .
EXPOSE 8080
CMD ["./main"]