# Purpose of Database

This directory contains the code and scripts necessary to interact with the Database, including database intitialisation, migration, and connections.

## Files and Structure

- **database.go** - Manages the **database connection** and provides **utility function (Execute)** to interact with the database combined with **query builder (query_builder.go)**. (CRUD Operations)

- **migrations/** - Contains database migration files that **define changes to the database schema**.

## Setup

- Ensure that the '.env' is **properly defined** with all **appropriate values**. Use '.env.example' as a template.

- Migration scripts are ran on setup **(checking if any database changes are needed to be added)**.
