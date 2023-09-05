# Star Sight

A complete example of a service (Stellar Service) created using Golang.

## Main goals

Stellar-Scope has three main goals:

1. ** Getting APOD data **: The service extracts metadata and a picture of the day from [APOD] on a daily basis (https://api.nasa.gov/).
   Then this data is securely stored in Minio, Postgresql.

2. ** HTTP API Server **: HTTP API, allows you to access all album recordings and extract certain recordings for a selected day.

3. **Docker Image **: For ease of deployment and isolation, Stellar-Scope is encapsulated in a Docker image.

## Getting to work

Follow these simple steps to run Stellar-Scope on your local computer:

1. ** Clone the repository **: Start by cloning this repository to your local computer.

2. **Setting environment variables **: Create a `.env` file in the root directory of the project and set the necessary environment variables.
   These variables should include your PostgreSQL and Minio database settings to ensure that the service works correctly.

3. **Start the Service **: Run the following command to start the project and start the service:
   This command initializes Docker containers and creates a database if it does not already exist.

4. **Get access to the service**: You can access it through your web browser at [http://localhost:8080 ](http://localhost:8080 ).

## Stopping the service

To stop the project, run the `make down` command. This will stop Docker containers.

## Service Registration

To view the service logs, run the Create Logs command.

## View the status of containers
To view the status of Docker containers, run the `make ps` command.