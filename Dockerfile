FROM listmonk/listmonk:latest

# Create required directories
RUN mkdir -p /listmonk/static/email-templates /listmonk/static/public

# Copy email templates
COPY ./static/email-templates/* /listmonk/static/email-templates/

# Copy public files
COPY ./static/public /listmonk/static/public/