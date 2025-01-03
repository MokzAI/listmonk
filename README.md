# Listmonk with Docker Compose for Coolify

## Introduction

This is the [docker compose file](https://github.com/coollabsio/coolify/blob/main/templates/compose/listmonk.yaml) for Listmonk from Coolify but with CORS enabled for localhost and mokz.ai domain so the Listmonk API could be accessed by my apps.

## Use in Coolify

Copy and paste the `docker-compose.yml` file into a new Docker Compose project in Coolify. Set the domain properly and deploy. 

## Custom static files

For customizing static files like opt-in email or landing page, see `/static` directory. For more info, read [this article](https://yasoob.me/posts/setting-up-listmonk-opensource-newsletter-mailing/#custom-static-files).

---

Running static files locally w/ hotreload.

1. Install `air`

```sh
go install github.com/air-verse/air@latest
export PATH=$PATH:$(go env GOPATH)/bin # Add this to your .zshrc or .bashrc
```

2. Run `air`

```sh
air
```

For example, if you wanted to preview the opt-in email, you would go to `http://localhost:8080/public/templates/opt-in.html`.