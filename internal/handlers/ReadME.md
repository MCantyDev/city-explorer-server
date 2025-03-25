# Purpose of Handlers

This directory contains the handler functions for the application. 

Each handler is resposible for:

- **Processing** an incoming HTTP request.

- **Interacting** with services or models.

- **Sending** a response back to the client.

## Files and Structure

- **User.go** - Handles user **signup**, profile **management**, and user profile **retrieval**. (Anything to do with the user that doesnt require authentication)

- **City.go** - Handles requests related to the city data **retrieval** and **management**. (Will check database for any existing city data before send request to external APIs)

- **Auth.go** - Handles user authentication processes, including **login**, **logout**, and **JWT management**.

## Usage

- Handlers are **registered in the routes** and are mapped to **specific HTTP methods** (GET, POST, etc).

- The handlers **interact** with the rest of the application to **appropriately handle** the request and **send** appropriate responses.
