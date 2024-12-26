terraform {
    required_providers {
        docker = {
            source = "kreuzwerker/docker"
            version = "~> 3.0"
        }
    }
}

provider "docker" {}

resource "docker_image" "go_docker_app" {
    name = "go-docker-app:latest"
    build {
        context = "."
        dockerfile = "Dockerfile"
    }
}

resource "docker_container" "go_docker_app" {
    name = "go-docker-app"
    image = docker_image.go_docker_app.name
    ports {
        internal = 8080
        external = 8080
    }
}