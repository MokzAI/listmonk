# Listmonk with Docker Compose

## Introduction

This is a simple docker-compose setup for [Listmonk](https://listmonk.app/), a standalone, self-hosted, newsletter and mailing list manager. It is fast, feature-rich, and packed into a single binary. It uses a PostgreSQL database to store data.

## Features

* Health check for all included services
* Sample configurations
* Local persitation of uploaded files and database

### Environments

* **Local testing** It includes a compose file for local development.
* **Production ready** It includes a compose file for production.
* **Production traefik** It includes a compose file for production with traefik as reverse proxy.

## Prerequisites

* Docker
* Docker Compose

## Getting Started

Copy the `.env.sample` file to `.env` and change the content to your needs.

```bash
cp .env.sample .env
```

### Locally

```bash
docker compose -f docker-compose.local.yml up -d
```

### Production

```bash
docker compose up -d
