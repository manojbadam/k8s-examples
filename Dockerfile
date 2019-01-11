FROM golang:alpine
RUN apk update && apk add curl --no-cache
RUN mkdir /app
ADD . /app/
WORKDIR /app 
RUN go build -o main .
CMD ["./main"]