FROM golang
WORKDIR /Users/geoffroy/Desktop/wowstatistician
COPY . .
RUN go build
EXPOSE 8080
CMD ["./wowstatistician" "serve"]