# Golang JWT Auth

A RESTful API service built in Go (Golang) that provides user authentication using JSON Web Tokens (JWT). This project demonstrates secure user login and registration, protected routes, and JWT token generation and validation in Golang.

## Features

- **User Registration**: Register a new user with username and password.
- **User Login**: Authenticate an existing user and receive a JWT.
- **JWT-Based Authentication**: Use JWT to protect routes.
- **Protected Routes**: Access to certain routes is restricted to authenticated users.
- **Token Validation**: Middleware for validating tokens on each request.

## Technologies Used

- **Golang**: Core language for building the API.
- **JWT (JSON Web Token)**: For secure token-based authentication.
- **Gorilla Mux**: Router for handling API routes.
- **bcrypt**: For hashing and verifying user passwords.

## Prerequisites

- [Go](https://golang.org/) installed (version 1.16+)
- A tool like [Postman](https://www.postman.com/) or [curl](https://curl.se/) for testing API requests.

## Installation

1. **Clone the repository**:
    ```bash
    git clone https://github.com/bhaveshs012/golang-jwt-auth.git
    cd golang-jwt-auth
    ```

2. **Install dependencies**:
    ```bash
    go mod tidy
    ```

3. **Set up environment variables**:
   Create a `.env` file in the root directory and add the following variables:
    ```bash
    MONGODB_URL=<your_mongodb_atlas_connection_string>
    PORT=4000
    SECRET_KEY=<your_jwt_secret>
    ```

4. **Run the server**:
    ```bash
    go run main.go
    ```

    The server will start on `http://localhost:4000`.
