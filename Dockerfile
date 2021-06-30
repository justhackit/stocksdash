#To Build (make sure GITHUB_PERSONAL_TOKEN is set): docker build --build-arg github_personal_access_token=$GITHUB_PERSONAL_TOKEN  -t stocksdash:v0.1 .
#To Push : docker tag stocksdash:v0.1 ajayedap/cloudlifter-images:stocksdash;docker push ajayedap/cloudlifter-images:stocksdash
#To Run on local: docker run -p 8090:8090 stocksdash:v0.1
FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN mkdir /build 
WORKDIR /build 

RUN go env -w GOPRIVATE="github.com/cloudlifter/*"

ARG github_personal_access_token
RUN git config \
  --global \
  url."https://justhackit:${github_personal_access_token}@github.com".insteadOf \
  "https://github.com"

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o stocksdash *.go

FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY --from=builder /build/stocksdash .
#COPY --from=builder /build/k8s-config.yaml .
#CMD ["./stocksdash","k8s-config.yaml"]