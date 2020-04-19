FROM golang:1.13-alpine3.10 as builder

RUN apk add --no-cache git live-media live-media-dev gcc g++ libc-dev gstreamer make cmake curl gstreamer-dev gst-plugins-base-dev ffmpeg-dev openssl-dev && \
    rm -rf /var/cache/apk/*
ENV CGO_ENABLED 1
ENV GOOS linux

ARG OPENCV_VERSION="4.3.0"
ENV OPENCV_VERSION $OPENCV_VERSION
RUN cd /tmp && \
    curl -Lo opencv.zip https://github.com/opencv/opencv/archive/${OPENCV_VERSION}.zip && \
    unzip -q opencv.zip && \
    curl -Lo opencv_contrib.zip https://github.com/opencv/opencv_contrib/archive/${OPENCV_VERSION}.zip && \
    unzip -q opencv_contrib.zip && \
    rm opencv.zip opencv_contrib.zip && \
    cd opencv-${OPENCV_VERSION} && \
    mkdir build && cd build && \
    cmake -D CMAKE_BUILD_TYPE=RELEASE \
    -D CMAKE_INSTALL_PREFIX=/tmp/base \
    -D OPENCV_EXTRA_MODULES_PATH=../../opencv_contrib-${OPENCV_VERSION}/modules \
    -D WITH_JASPER=OFF \
    -D WITH_QT=OFF \
    -D WITH_GTK=OFF \
    -D WITH_GSTREAMER=ON \
    -D BUILD_DOCS=OFF \
    -D BUILD_EXAMPLES=OFF \
    -D BUILD_TESTS=OFF \
    -D BUILD_PERF_TESTS=OFF \
    -D BUILD_opencv_java=NO \
    -D BUILD_opencv_python=NO \
    -D BUILD_opencv_python2=NO \
    -D BUILD_opencv_python3=NO \
    -D OPENCV_GENERATE_PKGCONFIG=ON .. && \
    make -j $(nproc --all) && \
    make preinstall && make install && \
    cd /tmp && rm -rf opencv*

ENV PKG_CONFIG_PATH /tmp/base/lib/pkgconfig:/tmp/base/lib64/pkgconfig
WORKDIR /build
ADD . .
RUN go build

# 64 bit platforms put it in lib64 which won't work
RUN [ -d /tmp/base/lib64 ] && mv /tmp/base/lib64 /tmp/base/lib || true

FROM alpine:3.10
RUN apk add --no-cache live-media libstdc++ ffmpeg-libs gstreamer gst-plugins-base gst-plugins-good gst-libav && \
    rm -rf /var/cache/apk/*
WORKDIR /opt/rts2p
COPY --from=builder /tmp/base/ /usr/
COPY --from=builder /build/rts2p /opt/rts2p/rts2p
COPY example.yaml /opt/rts2p/rts2p.yaml
CMD [ "/opt/rts2p/rts2p", "-c", "/opt/rts2p/rts2p.yaml" ]