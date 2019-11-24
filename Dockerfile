FROM golang AS build

COPY . /app
RUN cd /app && make

FROM scratch

COPY --from=build /app/url-parser /url-parser
ENTRYPOINT ["/url-parser"]
