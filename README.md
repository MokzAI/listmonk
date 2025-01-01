# Listmonk with Docker Compose for Coolify

## Introduction

This is the [docker compose file](https://github.com/coollabsio/coolify/blob/main/templates/compose/listmonk.yaml) for Listmonk from Coolify but with CORS enabled for localhost and mokz.ai domain so the Listmonk API could be accessed by my apps.

## Use in Coolify

Copy and paste the `docker-compose.yml` file into a new Docker Compose project in Coolify. Set the domain properly and deploy. 

## Custom static files

For customizing static files like opt-in email or landing page, see `/static` directory. For more info, read [this article](https://yasoob.me/posts/setting-up-listmonk-opensource-newsletter-mailing/#custom-static-files).

---

todo
- [ ] use of .env file - https://docs.docker.com/compose/how-tos/environment-variables/set-environment-variables/