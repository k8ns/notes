FROM scratch


RUN echo ${DB_USER}
RUN echo $${DB_USER}

COPY config.prod.yml /config.yml
COPY build/notes_api /

CMD ["/notes_api"]
