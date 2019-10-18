FROM scratch

COPY config/config.prod.yml /opt/notes_api/config/config.yml
COPY dest/notes_api /opt/notes_api/notes_api

WORKDIR /opt/notes_api

CMD ["/opt/notes_api/notes_api"]
