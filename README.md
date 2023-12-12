Lifetime and Daily Leaderboard for streaming backend application
Overview
This project is a representation of how the production level code should look like when you are building some microservice in Golang. Also for real time leaderboard applications like Gaming and YouTube, We can use sorted Sets data structure in Redis which is used in this project.

Features
Daily and Lifetime Views: Tracks views of videos on a daily and overall (lifetime) basis.
Scalable & Distributed: Utilizes Redis for caching and Consul for service discovery, making it scalable and easy to manage.
Onion Architecture: Organized in layers for modularity and ease of maintenance.
Tech Stack
Go (Golang): Main programming language used for development.
Redis: Used for caching video views and leaderboard scores.
Docker: Enables containerization for easy deployment and portability.
Consul: Configuration loading for distributed architecture
Architecture
The project is structured based on the Onion architecture, promoting clean and maintainable code. It's divided into layers:

Core: Contains the business logic and domain models.
Repositories: Handles data access and interactions with Redis.
Services: Implements service-level logic.
Transport: Manages HTTP endpoints.
Setup
Clone the Repository:
git clone https://github.com/Prthmesh6/Microservice_Gokit_Sample.git
Build & Run Using Docker:
docker compose up -d
Consul Configurations

{
      "key":"",
      "port":":",
      "address" : "",
      "poolSize":0,
      "username":"",
      "password": ""
}