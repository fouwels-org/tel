# SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
#
# SPDX-License-Identifier: MIT

FROM alpine:3.15.0 as build

# Install tools required for project
RUN apk add go git gcc make linux-headers cmake g++

ENV LIB61850_VERSION=1.5.0

# Build LIC61850
RUN git clone -c advice.detachedHead=false --depth 1 --branch v${LIB61850_VERSION} https://github.com/mz-automation/libiec61850.git
RUN cd libiec61850 && mkdir -p out
RUN cd libiec61850/out && cmake .. && make -j $(nproc) && make install
ENV GOBIN=/build/out

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ./go
RUN cd go && go install .

FROM alpine:3.15.0
COPY --from=build /build/out /build/out
COPY --from=build /usr/local/lib/ /usr/local/lib/
COPY --from=build /usr/local/include/ /usr/local/include/
RUN addgroup -S tel && adduser -S tel -G tel
USER tel

CMD ["/build/out/tel"]
