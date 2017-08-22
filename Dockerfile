FROM alpine:latest
#COPY ./tmp /tmp
EXPOSE 3000 80
COPY ./highloadcamp /highloadcamp
RUN /highloadcamp --validate
ENTRYPOINT ["/highloadcamp"]
CMD ["--port=80"]
