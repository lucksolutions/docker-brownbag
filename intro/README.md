# Docker Brownbag - Intro to Docker
![Docker Logo](https://www.docker.com/sites/default/files/moby.svg)

## What is Docker?
Docker is a tool that builds a "container" that contains everything needed to run your application. This includes
your code deployment, runtime, system tool and system libraries; anything that would be installed on your server. The goal of
the tool is to provide a configuration controlled, deployable artifact that can be guaranteed to run the same, no matter where
it is deployed.


### What are containers?
Containers run on a single machine and share the same OS kernel. Container images are constructed using a layered
filesystem and share common files. This makes the image disk space much more efficient, which means your image file sizes
quite small and quick to download.

### How are they different from VMs?
Virtual Machines contain an entire guest operating system in addition to your application binaries and libraries. 
![VM Architecture](https://www.docker.com/sites/default/files/WhatIsDocker_2_VMs_0-2_2.png) 

Containers
share the kernel with other contains and run as insolated process in the user space on the host OS. For more information on This
read up on LXE containers.
![Container Architecture](https://www.docker.com/sites/default/files/WhatIsDocker_3_Containers_2_0.png)

## Why would I want to use Docker?

### Developer Use Cases
* Docker provides an easy way for new team members to setup any software products or dependant services. For example your project
your project might require a MySQL database, Redis instance, or AEM Server.
* You can be sure that your code will run in any environment, because the runtime configuration and dependencies are controlled
in the container.
* Have multiple installations of application stacks for different projects. All while keeping your local machine clean.

### Production Use Cases
* Provides configuration management for server that can be easily managed in source control.
* Deploy the same container to any environment. Docker is supported in most major cloud providers; AWS, Azure, Digital Ocean, etc.
* Docker Trusted Registry will scan images for known security issues.
* Facilitates continuous delivery
* Platform supports auto-scaling, sevice discovery, load balancing, and rolling updates natively.


## Getting Started with Docker
Docker is supported on Linux, Mac and Windows machines. Download the platform at [https://www.docker.com/products/overview](https://www.docker.com/products/overview)

Docker architecture is comprised of the following components:
* Docker Engine - Interfaces with the host OS to run containers. 
* Docker Registry - Image repository. Primary registry is [Docker Hub](https://hub.docker.com/)
* Docker Compose - Used to define mult-container applications
* Docker Machine - Can be used to provision Docker on remote systems, or run multiple docker nodes locally.

### Docker Engine - Running your first container
After you have installed Docker, you can run your first container by with a single command:
`docker run -d -p 80:80 nginx:apline`
This will download the [official Nginx image](https://hub.docker.com/_/nginx/) from Docker Hub and start the container on your local Docker engine.

We can see all the running contains with the `docker ps` command. With the container running, we should be able
to bring up the Ngnix welcome page at [http://localhost](http://localhost). 

We can stop that container with the `docker stop <container name>`. This will leave the container stopped and it can be
restarted later with `docker start <container name>`. Since we don't need this container anymore, we will remove it using
`docker rm <container name>`

### Docker Engine - Building your first image
Instructions for building your image is defined in a `Dockerfile`. This is just a text document containing
a set of commands that a user could run on the command line to construct the image. The [Dockerfile Reference](https://docs.docker.com/engine/reference/builder/)
for a list of all avaiable commands you can use in your Dockerfile.

In our [Dockerfile](Dockerfile) we are extending the Nginx image and installing our static html content. We can then
build this image with the following command: `docker build -t custom-nginx .`

We can start our new container with: `docker run -d -p 80:80 custom-nginx`

For a more complex example we can look at this [Dockerfile](https://github.com/jasonluck/hybris-docker/blob/master/Dockerfile) I created for Hybris

## Important Things to Know About using Containers
* Start up order and timing is very important when your container relys on other services. You can use the _depends_on_ configuration to tell Docker which 
other services your container depends on, but that won't completely solve the problem. Because the container is considered "started" once the container CMD
is executed. If the service running in container takes time to completely become available, then the service maybe still not be ready by the time your other
container starts and tries to use the service.

* Data stored in containers is removed when the container is removed. If you want the data in the container
to persist past the life of the container, you need to make sure the data is stored in a [Volume](https://docs.docker.com/engine/tutorials/dockervolumes/).

* The exchange of secrets in a completely secure manner is a real challenge right now. Older practices stored secrets as `ENV` variables
within the container, which isn't completely secure. The latest version of Docker (v1.13) has added support for [secure secret storage](https://docs.docker.com/engine/swarm/secrets/) to the swarm engine.
Using a secrets management tool, such as [Vault](https://www.vaultproject.io) is another great option to completely secure secret data. How secrets are going
to be consumed by your application is something you need to carefully consider when writing containerized applications.

## Additional Resources
* [Brown Bag Presentation - https://github.com/jasonluck/docker-brownbag](https://github.com/jasonluck/docker-brownbag)
* [Docker Documentation - https://docs.docker.com/](https://docs.docker.com/)