FROM golang:1.14.2
ENV GO111MODULE "on"
ENV GOPROXY "https://goproxy.cn"
ADD . $GOPATH/src/github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/
WORKDIR $GOPATH/src/github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/
RUN make
EXPOSE 1203
CMD ["./main", "-c", "conf/config.yaml"]
