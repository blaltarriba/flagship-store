FROM iron/base

WORKDIR /app

COPY flagship-store /app/

EXPOSE 3080

ENTRYPOINT ["./flagship-store"]