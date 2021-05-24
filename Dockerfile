FROM golang:latest

ENV PROJECT github.com/marcochilese/negawordfixer
ENV GO111MODULE=off

ADD . $GOPATH/src/$PROJECT
RUN go get $PROJECT/...;
RUN cd $GOPATH/src/$PROJECT/ && go build -o main && mv main /main
RUN cp -r $GOPATH/src/$PROJECT/dictionary_data /
RUN chmod +x $GOPATH/src/$PROJECT/extract.sh
RUN chmod +x $GOPATH/src/$PROJECT/compress.sh
RUN cp -r $GOPATH/src/$PROJECT/extract.sh /
RUN cp -r $GOPATH/src/$PROJECT/compress.sh /

WORKDIR /
# Run the executable
ENTRYPOINT ["./main"]
