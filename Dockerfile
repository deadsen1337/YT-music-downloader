FROM golang:latest
ENV TZ=Europe/Moscow
RUN mkdir /WhiteBot
WORKDIR /WhiteBot
COPY . .

RUN apt-get -y update && apt-get install -y tzdata
RUN apt-get -y upgrade
RUN apt-get update &&  apt-get install -y ca-certificates --no-install-recommends && rm -rf /var/lib/apt/lists/*


CMD ["./bot"]