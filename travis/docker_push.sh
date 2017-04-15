#!/usr/bin/env bash

echo "docker_push.sh script start"

docker version

echo "docker repo: $DOCKER_REPO"
echo "docker commit tag: $DOCKER_COMMIT_TAG"
echo "docker build/push for TRAVIS_EVENT_TYPE=$TRAVIS_EVENT_TYPE"

docker login -e $DOCKER_EMAIL -u $DOCKER_USER -p $DOCKER_PASSWORD

if [ "$TRAVIS_EVENT_TYPE" == "pull_request" ]; then
  docker tag $DOCKER_REPO:$DOCKER_COMMIT_TAG $DOCKER_REPO:PR-$TRAVIS_PULL_REQUEST;
fi

if [ "$TRAVIS_EVENT_TYPE" == "push" ]; then
  docker tag $DOCKER_REPO:$DOCKER_COMMIT_TAG $DOCKER_REPO:$TRAVIS_BRANCH;
  if [ "$TRAVIS_BRANCH" == "master" ]; then
    docker tag $DOCKER_REPO:$DOCKER_COMMIT_TAG $DOCKER_REPO:latest;
  fi
fi

docker images
docker push $DOCKER_REPO
