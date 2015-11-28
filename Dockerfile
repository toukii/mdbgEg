FROM golang

WORKDIR /app/gopath/mdbg
ENV GOPATH /app/gopath

RUN git clone --depth 1 git://github.com/shaalx/mdbg.git .

RUN go get -u github.com/shaalx/mdbgEg/md
RUN mv $GOPATH/bin/md /bin/md

RUN go get -u github.com/shaalx/mdbgEg/rdr
RUN mv $GOPATH/bin/rdr /bin/rdr

RUN go get -u github.com/shaalx/mdbgEg
RUN mv $GOPATH/bin/mdbgEg /bin/mdbgEg

RUN mkdir -p /app/gopath/mdbg/MDFs

RUN rdr

EXPOSE 80

CMD ["/bin/mdbgEg"]


