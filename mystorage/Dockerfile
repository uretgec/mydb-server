FROM alpine:latest

LABEL service="mystorage"

RUN mkdir -p \
    /usr/local/mystorage \
    /usr/local/data/db/live/adminstore \
    /usr/local/data/db/live/clistore \
    /usr/local/data/db/live/livestore \
    /usr/local/data/db/live/moderatestore \
    /usr/local/data/db/live/pagestore

ARG ENV_FILE

COPY ${ENV_FILE} /usr/local/mystorage/.env
COPY --from=bakbibu-builder /usr/local/go/bin/mystorage /usr/local/mystorage/mystorage

WORKDIR /usr/local/mystorage

ENTRYPOINT [ "./mystorage", "--env-file", ".env" ]