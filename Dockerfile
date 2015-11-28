FROM golang

WORKDIR /app/gopath/mdbgEg
ENV GOPATH /app/gopath

RUN git clone --depth 1 git://github.com/shaalx/mdbgEg.git .

RUN go get -u github.com/shaalx/mdbgEg/md
RUN mv $GOPATH/bin/md /bin/md

RUN go get -u github.com/shaalx/mdbgEg/rdr
RUN mv $GOPATH/bin/rdr /bin/rdr

RUN go get -u github.com/shaalx/mdbgEg
RUN mv $GOPATH/bin/mdbgEg /bin/mdbgEg

RUN mkdir -p /app/gopath/mdbgEg/MDFs

RUN rdr

EXPOSE 80

CMD ["/bin/mdbgEg"]


