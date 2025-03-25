# Purpose of Models

This directory contains the data model representing the structure of the database tables. 

Models define the fields and relationships of the data objects used within the application.

## Files and Structure

- **User.go** - Represents the **user table** in the database, including fields such as:
  - **ID** - User's ID,
  - **Role** - User's role,
  - **FirstName** - User's first name,
  - **LastName** - User's last name,
  - **Email** - User's email,
  - **Password** - User's password.

- **City.go** - Represents the **city table** in the database, including fields such as:
  - **ID** - City's ID,
  - **Name** - City's name,
  - **Country** - Country where the city is located.

## Usage

- Models are used to **interact with the database** using the ORM. (GORM, from "github.com/go-gorm/gorm")

- Extend models to suit the **tables and relationships** defined in the database.
