FROM golang:1.13.6-alpine3.11 as build
RUN apk add --no-cache \
  git openssh-client \
  gcc musl-dev
ENV CGO_ENABLED 0

COPY . /usr/src/sharepoint2sync
WORKDIR /usr/src/sharepoint2sync
RUN go test ./... && go build -o sharepoint2sync -v ./cmd/sharepoint2sync.go


FROM debian:latest
ENV TZ Pacific/Noumea

ENTRYPOINT ["sharepoint2sync"]

COPY --from=build /usr/src/sharepoint2sync/sharepoint2sync /bin/