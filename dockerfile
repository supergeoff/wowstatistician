FROM ubuntu:latest
WORKDIR /usr/app/wowstatistician
COPY . .
CMD ./wowstatistician serve