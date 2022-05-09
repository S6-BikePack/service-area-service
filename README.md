<div align="center">

  <img src="assets/logo.png" alt="logo" width="200" height="auto" />
  <h1>BikePack - Service-Area-Service</h1>

  <p>
    Part of the S6 BikePack project.
  </p>


<!-- Badges -->

[![golangci-lint](https://github.com/S6-BikePack/service-area-service/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/S6-BikePack/service-area-service/actions/workflows/golangci-lint.yml)
[![Makefile CI](https://github.com/S6-BikePack/service-area-service/actions/workflows/run-tests.yml/badge.svg)](https://github.com/S6-BikePack/service-area-service/actions/workflows/run-tests.yml)
[![Docker](https://github.com/S6-BikePack/service-area-service/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/S6-BikePack/service-area-service/actions/workflows/docker-publish.yml)

<h4>
    <a href="https://github.com/S6-BikePack">Home</a>
  <span> · </span>
    <a href="https://github.com/S6-BikePack/service-area-service#-about-the-project">Documentation</a>
  </h4>
</div>

<br />

<!-- Table of Contents -->
# 📓 Table of Contents

- [About the Project](#-about-the-project)
    * [Architecture](#-architecture)
    * [Tech Stack](#%EF%B8%8F-tech-stack)
    * [Environment Variables](#-configuration)
    * [Messages](#-messages)
- [Getting Started](%EF%B8%8F-getting-started)
    * [Prerequisites](%EF%B8%8F-prerequisites)
    * [Running Tests](#-running-tests)
    * [Run Locally](#-run-locally)
    * [Deployment](#-deployment)
- [Usage](#-usage)




<!-- About the Project -->
## ⭐ About the Project

The Service-Area-Service is the service for the BikePack project that handles all service-areas in the system.
Using the system new service-areas can be created in which BikePack would operate.

<!-- Architecture -->
### 🏠 Architecture
For this service I have chosen a Hexagonal architecture. This keeps the service loosely coupled and thus flexible when having to change parts of the system.


<!-- TechStack -->
### 🛰️ Tech Stack
#### Language
  <ul>
    <li><a href="https://go.dev/">GoLang</a></li>
</ul>

#### Dependencies
  <ul>
    <li><a href="https://github.com/gin-gonic/gin">Gin</a><span> - Web framework</span></li>
    <li><a href="https://github.com/gin-gonic/gin">Amqp091-go</a><span> - Go AMQP 0.9.1 client</span></li>
    <li><a href="https://github.com/swaggo/swag">Swag</a><span> - Swagger documentation</span></li>
    <li><a href="https://gorm.io/index.html">GORM</a><span> - ORM library</span></li>
  </ul>

#### Database
  <ul>
    <li><a href="https://www.postgresql.org/">PostgreSQL</a></li>
</ul>

<!-- Env Variables -->
### 🔑 Configuration

This service can be configured using environment variables or a json file. The location of the file can be passed to the `config` environment variable.

The following configuration is available:

```json
{
    "server": {
      "service": "string",
      "port": "string",
      "description": "string"
    },
    "rabbitMQ": {
      "host": "string",
      "port": "int",
      "user": "string",
      "password": "string",
      "exchange": "string"
    },
    "database": {
      "host": "string",
      "port": "int",
      "user": "string",
      "password": "string",
      "database": "string",
      "debug": "bool"
    }
}
```

<!-- Data -->

##  🗃️ Data

This service stores the following data:

```json
{
  "ID": "int",
  "Identifier": "string",
  "Name": "string",
  "Area": "Polygon"
}
```

<!-- Messages -->
## 📨 Messages

### Publishing
The service publishes the following messages to the RabbitMQ server:

---
**service_area.create**

Published when a new service-area is created in the system.
Sends the newly created service-area in the  body.

```json
{
  "ID": "int",
  "Identifier": "string",
  "Name": "string",
  "Area": "Polygon"
}
```

---
**service_area.update**

Published when a service-area is updated in the system.
Sends the updated service-area in the body.

```json
{
  "ID": "int",
  "Identifier": "string",
  "Name": "string",
  "Area": "Polygon"
}
```

<!-- Getting Started -->
## 	🛠️ Getting Started

<!-- Prerequisites -->
### ‼️ Prerequisites

Building the project requires Go 1.18.

This project requires a PostgreSQL compatible database with a database named `service-area` and a RabbitMQ server.
The easiest way to setup the project is to use the Docker-Compose file from the infrastructure repository.

<!-- Running Tests -->
### 🧪 Running Tests

The tests in the project can easily be run using make and the `make run-tests` command. This will start the required docker containers and run all tests in the project.

<!-- Run Locally -->
### 🏃 Run Locally

Clone the project

```bash
  git clone https://github.com/S6-BikePack/service-area-service
```

Go to the project directory

```bash
  cd service-area-service
```

Run the project (Rest)

```bash
  go run cmd/rest/main.go
```


<!-- Deployment -->
### 🚀 Deployment

To build this project run (Rest)

```bash
  go build cmd/rest/main.go
```


<!-- Usage -->
## 👀 Usage

### REST
Once the service is running you can find its swagger documentation with all the endpoints at `/swagger`