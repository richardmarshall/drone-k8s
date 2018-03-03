# drone-k8s

[![Build Status](https://travis-ci.org/richardmarshall/drone-k8s.svg?branch=master)](https://travis-ci.org/richardmarshall/drone-k8s)

Drone k8s plugin

This plugin can be used to apply manifests to a kubernetes cluster.

# Example configuration

```yaml
pipeline:
  kubernetes:
    image: richardmarshall/drone-k8s
    server: https://k8s.example.com
    secrets: [ k8s_token ]
```

# Template 

# Secret Reference

`k8s_token`  
authenticates with this token

`k8s_ca`  
base64 encoded ca public certificate

# Parameter Reference

`manifest=deployment.yaml`  
file containing kubernetes objects to be applied

`server`  
kubernetes API endpoint

`token`
authenticates with this token

`ca`  
base64 encoded ca public certificate

`namespace`  
apply this namespace to objects without one explicitly set

`prune=false`  
prune objects matching `selector` that are not in the manifest

`selector`
label selector for finding objects to prune