FROM alpine:latest
COPY ./highloadcamp /highloadcamp
RUN /highloadcamp --validate
EXPOSE 80
ENTRYPOINT ["/highloadcamp"]
CMD ["--port=80"]
