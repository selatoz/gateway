# Go Gin Project Template

This is a bare Go Gin project template with a pre-defined directory structure that separates the app's business logic, HTTP request handlers, configuration logic, database abstraction layer, and routes abstraction layer.

## What is Gin?

Gin is a web framework written in Go (Golang) that's designed to be lightweight and efficient. It's built on top of the net/http package, and it provides a simple and flexible API for building web applications.

To learn more about Gin, you can check out the official documentation at [https://gin-gonic.com/](https://gin-gonic.com/).

## Directory Structure

The project directory structure is as follows:

- `./api` and `./api/auth` - where the HTTP request handlers are located
- `./config` - where the configuration logic is located
- `./database` - where the database abstraction layer for db setup is located
- `./routes` - where the routes abstraction layer for the app is located
- `./internal` with `./internal/svc` and `./internal/repo` - where the app's business logic is located (svc stands for service, and repo for repository)

By separating the different parts of the application into separate directories, it becomes easier to manage and maintain the codebase. The `./api` directory contains the HTTP request handlers, which are responsible for handling incoming requests and generating responses. The `./config` directory contains the configuration logic, which handles reading configuration files and setting up environment variables. The `./database` directory contains the database abstraction layer, which handles database connections and setup. The `./routes` directory contains the routes abstraction layer, which defines the API routes for the app. Finally, the `./internal` directory contains the app's business logic, which is further split into `./internal/svc` for the service layer and `./internal/repo` for the repository layer.

By following this directory structure, you can create a well-organized and maintainable Go Gin project that's easy to work with and extend.

## Database Management

This project template utilizes `gorm.DB`, a popular Object-Relational Mapping library for Go, for managing database connections and interactions. Gorm provides a convenient and expressive API for working with databases, allowing for efficient and easy database management.

To learn more about `gorm.DB` and how to use it, you can check out the official documentation at [https://gorm.io/](https://gorm.io/).

