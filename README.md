# About the service

A complete example of a service (Stellar Service) created using Golang.

## Main goals

Stellar-Scope has three main goals:

1. **Getting APOD data**: The service extracts metadata and a picture of the day from [APOD](https://api.nasa.gov/) on a daily basis.
   Then this data is securely stored in Minio, Postgresql.

2. **HTTP API Server**: HTTP API, allows you to access all album recordings and extract certain recordings for a selected day.

3. **Docker Image**: For ease of deployment and isolation, Stellar-Scope is encapsulated in a Docker image.

## Getting to work

Follow these simple steps to run Stellar-Scope on your local computer:

1. **Clone the repository**: Start by cloning this repository to your local computer.

2. **Start the Service**: Run the following command `make up` to start the project.
   This command initializes Docker containers and creates a database if it does not already exist.

## Stopping the service

To stop the project, run the `make down` command. This will stop Docker containers.

## Logging

To view the service logs, run the`make logs` command.

## View the status of containers
To view the status of Docker containers, run the `make ps` command.