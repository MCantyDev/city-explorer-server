# Purpose of Middleware

This directory contains the middleware functions used by the **Gin web server**.

Middleware functions are executed during the **request-response cycle** and can **modify** the **request** or the **response**.

## Files and Structure

- **Auth.go** - Middleware for **handling** the JWT **authentication and authorisation checks**. (executed on a per request basis, uses the Auth service)

- **Logging.go** - Middleware for **centralised logging** of request details. (e.g. Method[GET, POST, etc.], URL, Response Time)

- **ErrorHandler.go** - Middleware for handling and formatting errors in a **consistent way** throughout the application.

## Usage

- Each middleware is **defined in the router setup** (main.go).
