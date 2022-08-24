ARG BUILD_IMAGE=golang:1.19-alpine

FROM $BUILD_IMAGE as build
WORKDIR /workspace
ENV GOPROXY=https://goproxy.cn
COPY go.mod go.sum /workspace/
RUN go mod download
COPY cmd /workspace/cmd
COPY dist /workspace/dist
COPY pkg /workspace/pkg
ARG GOFLAGS
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${GOFLAGS} -o vmdashboard cmd/dashboard/main.go

FROM alpine:3.14
COPY --from=build /workspace/vmdashboard /usr/local/bin/vmdashboard
ENTRYPOINT [ "/usr/local/bin/vmdashboard" ]