# gRPC_SSO

gRPC_SSO is a project that provides a gRPC service for Single Sign-On (SSO) authentication. It is implemented using the Go programming language and relies on PostgreSQL as its database. For deployment convenience, Docker and Docker Compose files are included.

## Features

- Provides a gRPC service for SSO authentication
- Uses PostgreSQL for database management
- Simplifies deployment with Docker

## Before running the application, ensure you have the following installed:

- Docker

## Installation and Setup

1. **Clone the repository:**

    ```bash
    git clone https://github.com/yourusername/gRPC_SSO.git
    ```

2. **Navigate to the project directory:**

     ```bash
    cd gRPC_SSO
    ```
3. **Modify the configuration in config/local.yaml according to your requirements**
4. **Run the application using Docker:**
   ```bash
   docker-compose up
   ```
5. **Access the gRPC service:**
   The gRPC service will be available at localhost:your_port by default.

## Configuration

Configuration options can be found in config/local.yaml and docker-compose file.

## API Reference

For details on the gRPC API provided by this service, refer to the protos/proto/sso.proto file.
