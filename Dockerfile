FROM golang:1.19.5-alpine3.17

ARG VIPS_VERSION="8.14.1"

RUN wget https://github.com/libvips/libvips/releases/download/v${VIPS_VERSION}/vips-${VIPS_VERSION}.tar.xz
RUN apk update && apk add meson build-base pkgconfig glib-dev gobject-introspection-dev libxml2-dev expat-dev jpeg-dev libwebp-dev libpng-dev

RUN tar -xf vips-${VIPS_VERSION}.tar.xz
WORKDIR /go/vips-${VIPS_VERSION}
RUN meson build
WORKDIR /go/vips-${VIPS_VERSION}/build
RUN meson compile && \
    meson test && \
    meson install

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/app ./cmd/redbud-way-api-server/main.go
EXPOSE 80

CMD ["app", "--host", "0.0.0.0", "--port", "80"]