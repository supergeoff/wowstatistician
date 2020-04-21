FROM golang
WORKDIR /usr/local/go/src/wowstatistician
COPY . .
RUN go build
EXPOSE 8080
CMD ./wowstatistician serve