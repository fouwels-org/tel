# SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
#
# SPDX-License-Identifier: MIT

FROM alpine:3.14.0 as build

# Install tools required for project
RUN apk add go

ENV GOBIN=/build/out

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ./go
RUN cd go && go install .

FROM alpine:3.14.0
COPY --from=build /build/out /build/out

RUN addgroup -S tel && adduser -S tel -G tel
USER tel

CMD ["/build/out/tel"]
