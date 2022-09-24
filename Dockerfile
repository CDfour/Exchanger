FROM golang

WORKDIR C:\go\apps\exchanger

COPY . .

RUN go build cmd/main.go

EXPOSE 6080

CMD ["main"]