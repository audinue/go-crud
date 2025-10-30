FROM golang:bookworm

RUN go install golang.org/x/tools/gopls@latest \
  && go install honnef.co/go/tools/cmd/staticcheck@latest

RUN go install github.com/audinue/gomon@latest
