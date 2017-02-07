# Docker Brownbag - Docker Compose
Docker Compose is a tool for defining and running multi-container applications. Inside a single file your define
all your application components and the application services they rely on. Then using a single command, you 
can create, start and stop all those services.

The compose file becomes the definition of your applications stack and its configuration. This file lives
in source control along side your application.

## Defining your application stack
A docker-compose.yml file looks like:
```
version: '2'

services:
    app-db:
        image: mysql
        environment: 
            MYSQL_RANDOM_ROOT_PASSWORD: 'yes'
            MYSQL_DATABASE: 'app'
            MYSQL_USER: 'appuser'
            MYSQL_PASSWORD: 'apppassword'
        networks:
            - backend

    app-db-schema:
        build:
            context: ./db
        image: dockerbrownbag-schema
        environment:
            DB_DATABASE: 'app' 
            DB_USER: 'appuser'
            DB_PASSWORD: 'apppassword'
        links:
            - app-db
        depends_on:
            - app-db
        networks:
            - backend

networks:
    backend:
        driver: bridge
```
Its contains a `service` entry for each container that will be needed. These are either container images from
a repository, or containers built from a local Dockerfile definition. 

It also contains a set of network definitions that can we used to segregate container communications. 

### Common Commands
* `docker-compose up` - Starts all services defined in your compose file. Use the `-d` option to run them in the background.
* `docker-compose down` - Stops and removes all the containers, networks and volumes associated with your compose file.
* `docker-compose ps` - Shows all running containers associated to your compose file.
* `docker-compose logs <service name>` - Shows the logs for the given service.

## Using Compose in Multiple Environments
This is all great for your local dev machine, but we probably need different environment configurations
for CI, Testing or Production environments.
