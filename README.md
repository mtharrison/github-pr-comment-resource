# github-pr-comment-resource

![master](https://github.com/mtharrison/github-pr-comment-resource/workflows/Master/badge.svg?branch=master) ![Release](https://github.com/mtharrison/github-pr-comment-resource/workflows/Release/badge.svg) ![GitHub release](https://img.shields.io/github/v/release/mtharrison/github-pr-comment-resource)

A resource type for [Concourse CI](https://concourse-ci.org/) to trigger builds from Github PR comments. Also allows parameters to be provided in comments.
 
![image](https://cldup.com/beeBL0NNQ3.png)
 
 ---
 
## Usage
 
 To use register the resource type using the public [Docker image](https://hub.docker.com/repository/docker/mtharrison/github-pr-comment-resource) `mtharrison/github-pr-comment-resource`.

 ```yaml
 resource_types:
  - name: github-pr-comment-resource
    type: docker-image
    source:
      repository: mtharrison/github-pr-comment-resource
      tag: v0.2.0
 ```
 Then create a new resource using this resource type. You'll need to provide the `repository`, `access_token`, optionally a `v3_endpoint` if you're using Github Enterprise, otherwise this will default to the public Github API.
 
 Optionally you can provide a `regex` (valid [Go regex](https://golang.org/pkg/regexp/) only). Comments that match this regex only will become new versions. If the regex contains capture groups they will be provided by the resource too.

 ```yaml
resources:
  - name: deployment-trigger
    type: github-pr-comment-resource
    icon: github
    source:
      repository: golang/go
      access_token: '[...]'
      v3_endpoint: '[...]'
      regex: '^deploy ([a-zA-Z0-9_.-]+) to ([a-zA-Z0-9_.-]+) please$'
 ```
 Finally you can use the resource as an input/trigger to any jobs.
 ```yaml
jobs:
  - name: deployment-test
    plan:
      - get: deployment-trigger
        trigger: true
      - task: echo
        config:
          image_resource:
            type: registry-image
            source:
              repository: alpine
          inputs:
            - name: deployment-trigger
          platform: linux
          run:
            path: /bin/sh
            args:
              - '-c'
              - |
                apk add jq &> /dev/null
                cat deployment-trigger/comment.json | jq
 ```
 In the resource directory a `comment.json` file will be written containing the resource data. See screenshot below:
 
 ![screenshot](https://cldup.com/ZyLNgJX85r.png)
