# Listmonk with Docker Compose

## Introduction

This is a simple docker-compose setup for [Listmonk](https://listmonk.app/), a standalone, self-hosted, newsletter and mailing list manager. It is fast, feature-rich, and packed into a single binary. It uses a PostgreSQL database to store data.

## Features

* Health check for all included services
* Sample configurations
* Local persitation of uploaded files and database

## Prerequisites

* Docker
* Docker Compose

## Getting Started

Copy the `.env.sample` file to `.env` and change the content to your needs.

```bash
cp .env.sample .env
```

## Run

```bash
docker compose up -d
```
