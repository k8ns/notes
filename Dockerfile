FROM scratch

COPY config.prod.yml /config.yml
COPY build/notes_api /

CMD ["/notes_api"]
