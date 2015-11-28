FROM golang

WORKDIR /app/gopath/mdbgEg
ENV GOPATH /app/gopath
#RUN export PATH=$PATH:$GOPATH/bin

RUN git clone --depth 1 git://github.com/shaalx/mdbgEg.git .

RUN go get -u github.com/shaalx/mdbgEg
RUN mv $GOPATH/bin/mdbgEg /bin/mdbgEg

RUN mkdir -p /app/gopath/mdbgEg/MDFs

EXPOSE 80

CMD ["/bin/mdbgEg"]


