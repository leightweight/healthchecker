FROM golang:1.20 as build

ARG VERSION
ARG DATE
ENV CGO_ENABLED=0

WORKDIR /go/src/healthchecker
COPY . .

RUN make VERSION="${VERSION}" DATE="${DATE}" cli


FROM gcr.io/distroless/base-debian11:nonroot

ENV HEALTHCHECKER_NOPRETTY=true

USER nonroot

COPY --from=build --chown=nonroot /go/src/healthchecker/healthchecker /usr/local/bin/

ENTRYPOINT ["healthchecker"]
CMD ["wait"]
