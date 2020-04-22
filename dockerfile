FROM ubuntu:latest
WORKDIR /usr/app/wowstatistician
COPY . .
RUN ls -als
CMD ./wowstatistician serve