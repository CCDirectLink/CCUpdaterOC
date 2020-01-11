FROM golang AS builder

RUN apt-get update && \
    apt-get install -y libsdl2-dev mingw-w64 clang llvm-dev libxml2-dev uuid-dev libssl-dev bash patch make tar xz-utils bzip2 gzip sed cpio libbz2-dev cmake gcc g++ zlib1g-dev libmpc-dev libmpfr-dev libgmp-dev && \
    apt-get clean

# Windows crosscompiler
RUN wget "http://libsdl.org/release/SDL2-devel-2.0.10-mingw.tar.gz" -O /SDL2-devel-2.0.10-mingw.tar.gz && \
    tar -xzf /SDL2-devel-2.0.10-mingw.tar.gz -C / &&\
    cp -r /SDL2-2.0.10/x86_64-w64-mingw32 /usr

# Mac crosscompiler

# ENV OSX_SDK MacOSX10.13.sdk
# ENV DARWIN_VERSION=15
# ENV OSX_VERSION_MIN=10.10
# RUN set -x \
# 	&& export OSXCROSS_PATH="/osxcross" \
# 	&& git clone https://github.com/tpoechtrager/osxcross.git $OSXCROSS_PATH \
# 	&& cd $OSXCROSS_PATH \
#   && wget "https://github.com/phracker/MacOSX-SDKs/releases/download/10.13/MacOSX10.13.sdk.tar.xz" -O "/osxcross/tarballs/MacOSX10.13.sdk.tar.xz" \
#	&& UNATTENDED=yes ${OSXCROSS_PATH}/build.sh
# ENV PATH /osxcross/target/bin:$PATH

WORKDIR /ccmu/
COPY . .

RUN go mod download
RUN CGO_ENABLED=1 CC=gcc GOOS=linux GOARCH=amd64 go build -v -tags static  -ldflags "-s -w" -o CCUpdaterUI .
RUN CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -v -tags static -ldflags "-s -w -H=windowsgui" -o CCUpdaterUI.exe .
# RUN CGO_ENABLED=1 CC=o64-clang CXX="o64-clang++" MACOSX_DEPLOYMENT_TARGET=10.13 GOOS=darwin GOARCH=amd64 go build -tags static -ldflags "-s -w" -o CCUpdaterUI_mac .

FROM halverneus/static-file-server

COPY --from=builder /ccmu/CCUpdaterUI /web/
COPY --from=builder /ccmu/CCUpdaterUI.exe /web/
# COPY --from=builder /ccmu/CCUpdaterUI_mac /web/