FROM ubuntu

RUN env
RUN echo ${DB_USER}
RUN echo $${DB_USER}

FROM scratch

COPY config.prod.yml /config.yml
COPY build/notes_api /

CMD ["/notes_api"]
