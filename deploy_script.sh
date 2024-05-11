#!/bin/bash

DOCKER_IMAGE="$DOCKERHUB_USERNAME/$SERVICE_NAME"
DOCKERHUB_IMAGE="$DOCKERHUB_USERNAME/$SERVICE_NAME:latest"
echo "$DOCKERHUB_IMAGE"

# Pull latest image from Docker Hub
echo "$DOCKERHUB_TOKEN" | docker login -u "$DOCKERHUB_USERNAME" --password-stdin
docker pull $DOCKERHUB_IMAGE
docker rmi $(docker images | grep $DOCKER_IMAGE | grep "<none>" | awk '{print $3}')

# Stop and remove the old container if it exists
if docker ps -a | grep -q 'go-boilerplate'; then
    docker stop go-boilerplate
    docker rm go-boilerplate
fi

# Run new container from the latest image
docker run --network app-network --name go-boilerplate -d -v $VOLUME_LOG:/app/log -p $SERVICE_PORT:3001 -v $VOLUME_CONFIG:/app/resources $DOCKERHUB_IMAGE
