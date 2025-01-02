FROM listmonk/listmonk:latest
COPY ./static/* /listmonk/static/
RUN ls -la /listmonk/static/