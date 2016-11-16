FROM alpine:3.4
RUN apk --update add ca-certificates
COPY ./build/linux-static/glio /bin/glio
COPY ./ui /ui
COPY ./com /com
COPY ./run /config
WORKDIR /
ENV LOCAL false
ENTRYPOINT ["/bin/glio"]
CMD ["/config/heroku.toml"]
