FROM debian
COPY ./matrix-app app
RUN chmod 777 app
ENTRYPOINT /app