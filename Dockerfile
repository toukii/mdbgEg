FROM golang

WORKDIR /app/gopath/mdbgEg
ENV GOPATH /app/gopath
#RUN export PATH=$PATH:$GOPATH/bin

RUN git clone --depth 1 git://github.com/toukii/mdbgEg.git .

RUN go get -u github.com/toukii/mdbgEg
RUN mv $GOPATH/bin/mdbgEg /bin/mdbgEg

RUN mkdir -p /app/gopath/mdbgEg/MDFs

EXPOSE 80

CMD ["/bin/mdbgEg"]


