FROM google/golang

MAINTAINER Lyan Hung <lyanchih@gmail.com>

RUN apt-get update && apt-get install -y pkg-config libxml2-dev && rm -rf /var/lib/apt/lists/*

WORKDIR /home/go

ADD . /home/go

RUN go get github.com/robfig/cron github.com/lyanchih/goamf code.google.com/p/go-uuid/uuid github.com/moovweb/gokogiri && go build -o /usr/bin/freeway

VOLUME ["/home/go/data"]

CMD ["freeway"]
