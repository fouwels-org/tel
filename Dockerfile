# SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
#
# SPDX-License-Identifier: MIT

FROM alpine:3.15.0 as build

# Install tools required for project
RUN apk add go git gcc make linux-headers

ENV LIB61850_VERSION=1.5.0

# Build LIC61850
RUN git clone --depth 1 --branch v${LIB61850_VERSION} https://github.com/mz-automation/libiec61850.git
RUN cd libiec61850 && make -j $(nproc) && make install
ENV GOBIN=/build/out

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ./go
RUN cd go && go install .

FROM alpine:3.15.0
COPY --from=build /build/out /build/out

RUN addgroup -S tel && adduser -S tel -G tel
USER tel

CMD ["/build/out/tel"]
