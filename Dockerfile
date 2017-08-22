FROM alpine:latest
COPY ./highloadcamp /highloadcamp
COPY ./tmp /tmp
RUN /highloadcamp --validate
EXPOSE 3000 80
ENTRYPOINT ["/highloadcamp"]
CMD ["--port=80"]
