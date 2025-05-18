# City Explorer Server

A backend server built with Go to provide location-based data for the City Explorer application.

---

## Features

- User authentication (Signup, Login, Logout)
- Secure API key handling for external APIs
- Location-based data retrieval
- JWT-based authentication for API requests
- Session Refreshing with Refresh Tokens
- Database integration
- Cookies

## Tech Stack

- **Go (Golang)** - Backend
- **Gin** - Web Framework
- **Gorm with MySQL Driver** - MySQL Database 
- **React** - Frontend (separate repository)
- **JWT (JSON Web Tokens)** - Authentication
- **External APIs** - (e.g., weather, maps, points of interest)

---

## Prerequisites

- Go installed (`>=1.24.1`)

## Setup Instructions

### Installation

1. Clone the repository:

```sh
git clone https://github.com/your-username/city-explorer-server.git
cd city-explorer-server
```

---

2. Create .env in project root based on .env.example

---

3. Build the Application:

```sh
go build -o server.exe cmd/main.go
```

---

---

4. Run the Application:

```sh
.\server
```

---

# API Endpoints

This document describes the API endpoints for the application, organised by access level.

---

## Public Endpoints

These endpoints are accessible without authentication.

| Method | Endpoint     | Description                |
|--------|--------------|----------------------------|
| POST   | `/login`     | User login                 |
| POST   | `/sign-up`   | Create a new user account  |

---

## Authenticated Endpoints

These endpoints require a valid session and are protected by `SessionAuthMiddleware`.

| Method | Endpoint                  | Description                                   |
|--------|---------------------------|-----------------------------------------------|
| GET    | `/auth/profile`           | Get the authenticated user's profile          |
| GET    | `/auth/logout`            | Log out the user and clear session cookies    |
| GET    | `/auth/get-country`       | Retrieve a list of supported countries        |
| GET    | `/auth/get-cities`        | Get cities associated with a input query      |
| GET    | `/auth/get-city-weather`  | Get current weather data for a specific city  |
| GET    | `/auth/get-city-sights`   | Get tourist sights available in a city        |
| GET    | `/auth/get-city-poi`      | Get points of interest (POIs) for a city      |

---

## Admin Endpoints

These endpoints require both authentication and admin privileges, enforced by `AdminMiddleware`.

### GET Requests

| Method | Endpoint                     | Description                                 |
|--------|------------------------------|---------------------------------------------|
| GET    | `/admin/get-users`           | List all users in the system                |
| GET    | `/admin/get-countries`       | List all countries in the database          |
| GET    | `/admin/get-city-weather`    | Retrieve all city weather records           |
| GET    | `/admin/get-city-sights`     | Retrieve all city sights records            |
| GET    | `/admin/get-city-pois`       | Retrieve all city POIs                      |

### POST Requests

| Method | Endpoint             | Description                |
|--------|----------------------|----------------------------|
| POST   | `/admin/add-user`    | Create a new user account  |

### PATCH Requests

| Method | Endpoint                          | Description                              |
|--------|-----------------------------------|------------------------------------------|
| PATCH  | `/admin/edit-user`                | Update an existing user's information    |
| PATCH  | `/admin/refresh-country`          | Refresh country dataset                  |
| PATCH  | `/admin/refresh-city-weather`     | Refresh city weather data                |
| PATCH  | `/admin/refresh-city-sights`      | Refresh city sights data                 |
| PATCH  | `/admin/refresh-city-poi`         | Refresh city points of interest (POIs)   |

### DELETE Requests

| Method | Endpoint                          | Description                              |
|--------|-----------------------------------|------------------------------------------|
| DELETE | `/admin/delete-user`              | Remove a user from the system            |
| DELETE | `/admin/delete-country`           | Remove a country from the dataset        |
| DELETE | `/admin/delete-city-weather`      | Delete weather data for a city           |
| DELETE | `/admin/delete-city-sights`       | Delete sights data for a city            |
| DELETE | `/admin/delete-city-poi`          | Delete points of interest for a city     |


## License

This project is licensed under the MIT License.

---
