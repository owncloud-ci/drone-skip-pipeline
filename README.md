# drone-skip-pipeline

[![Build Status](https://img.shields.io/drone/build/owncloud-ci/drone-skip-pipeline?logo=drone&server=https%3A%2F%2Fdrone.owncloud.com)](https://drone.owncloud.com/owncloud-ci/drone-skip-pipeline)
[![Docker Hub](https://img.shields.io/docker/v/owncloudci/drone-skip-pipeline?logo=docker&label=dockerhub&sort=semver&logoColor=white)](https://hub.docker.com/r/owncloudci/drone-skip-pipeline)
[![GitHub contributors](https://img.shields.io/github/contributors/owncloud-ci/drone-skip-pipeline)](https://github.com/owncloud-ci/drone-skip-pipeline/graphs/contributors)
[![Source: GitHub](https://img.shields.io/badge/source-github-blue.svg?logo=github&logoColor=white)](https://github.com/owncloud-ci/drone-skip-pipeline)
[![License: Apache-2.0](https://img.shields.io/github/license/owncloud-ci/drone-skip-pipeline)](https://github.com/owncloud-ci/drone-skip-pipeline/blob/main/LICENSE)

Drone plugin to skip pipelines based on changed files.

## Build

Build the binary with the following command:

```console
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/linux/amd64/drone-skip-pipeline
```

## Docker

Build the Docker image with the following command:

```console
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/Dockerfile.linux.amd64 --tag plugins/skip-pipeline .
```

## Usage

```console
docker run --rm \
  -e PLUGIN_DISALLOW_SKIP_CHANGED="^.drone.star$" \
  -e PLUGIN_ALLOW_SKIP_CHANGED="^docs/.*,^changelog/.*" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/skip-pipeline
```

## Releases

Please create and commit a changelog for the next tag first:

```Shell
git-chglog -o CHANGELOG.md --next-tag v2.10.3 v2.10.3
git add CHANGELOG.md; git commit -m "[skip ci] update changelog"; git push
```

Afterwards create and push the new tag to trigger the CI release process:

```Shell
git tag v2.10.3
git push origin v2.10.3
```

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](https://github.com/owncloud-ci/drone-skip-pipeline/blob/main/LICENSE) file for details.

## Copyright

```Text
Copyright (c) 2021 ownCloud GmbH
```
