Susbsription Service API
A service for create, read, update and delete subscription on your platform

--Instalation
Clone the repository:
git clone https://github.com/brokkol-12/subscription-service.git
cd subscription-service

--Environment Setup
Before launching the application, create a .env file in the root directory.
You can use the example file:
cp .env.example .env

Then edit the following variables:

DB_USER=your_db_user
DB_PASSWORD=your_db_password


Make sure the credentials match the PostgreSQL configuration.

Run with Docker (Recommended)

Build and start the containers:

docker compose up --build

After startup, the service will be available at:
http://localhost:8081

--Without docker
You need install on your PC:
GO 1.21+
PostgreSQL

after in your CMD enter the command in dirrectory with API
go mod tidy

and start a programm:
go run cmd/main.go

--Swagger
After start you get access to SWAGGER
Open your browser and in search panel enter: http://localhost:8081/swagger/index.html
Here you got any options for manipulate with service.
API Documentation

Swagger UI provides:
Create subscription
Get subscription by ID
Get all subscriptions
Update subscription
Delete subscription
Get a total sum on all subscription

--Tech Stack
Go 1.21+
PostgreSQL
Docker & Docker Compose
Swagger

--Notes
Make sure PostgreSQL is running before starting the app without Docker.
The .env file is required for database connection.
Do not commit .env to version control.
