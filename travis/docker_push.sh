#!/usr/bin/env bash
# TODO: handle pull request
echo "docker_push.sh script start"
docker version

echo "DOCKER_REPO: $DOCKER_REPO"
echo "DOCKER_COMMIT_TAG: $DOCKER_COMMIT_TAG"
echo "TRAVIS_EVENT_TYPE=$TRAVIS_EVENT_TYPE"
echo "BRANCH_REGEX: $BRANCH_REGEX"
echo "TRAVIS_BRANCH: $TRAVIS_BRANCH"
echo "TRAVIS_EVENT_TYPE: $TRAVIS_EVENT_TYPE"

docker login -e $DOCKER_EMAIL -u $DOCKER_USER -p $DOCKER_PASSWORD

if [ "$TRAVIS_EVENT_TYPE" == "pull_request" ]; then
  docker tag $GOCLOUD_PROJECT_NAME:$DOCKER_COMMIT_TAG $DOCKER_REPO:PR-$TRAVIS_PULL_REQUEST;
fi

if [ "$TRAVIS_EVENT_TYPE" == "push" ]; then
  docker tag $GOCLOUD_PROJECT_NAME:$DOCKER_COMMIT_TAG $DOCKER_REPO:$TRAVIS_BRANCH;
  if [ "$TRAVIS_BRANCH" == "master" ]; then
    docker tag $GOCLOUD_PROJECT_NAME:$DOCKER_COMMIT_TAG $DOCKER_REPO:latest;
  fi
fi

docker images
docker push $DOCKER_REPO
