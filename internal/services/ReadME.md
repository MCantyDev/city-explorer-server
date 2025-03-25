# Purpose of Services

This directory contains the functions used by the handlers. 

Services perform the complex operations (e.g. generating JWT tokens), often required to interact with different models and external APIs.

## Files and Structure

- **AuthService.go** - Contains the functions for **user authentication**. (e.g. generating JWT tokens, validating user credentials, etc.)

- **CityService.go** - Contains the function for **city data operations**. (e.g. checking database for existing data, fetching city data from external APIs)

## Usage

- Services are executed by handlers to **perform application logic** and return data or perform actions.

- Each service should be kept focused on **one area of functionality** (single responsibilty - The Auth service should NOT handle City data)
