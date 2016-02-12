# bosh-lite-ami-resource

A Concourse resource for discovering the Amazon Machine Image (AMI) for the cloudfoundry/bosh-lite Vagrant box

## Source config
```
---
resources:
- name: bosh-lite-version
  type: bosh-lite-ami
  source:
    region: us-west-2
```


## Implemented actions

### check
Returns the version of the latest box

### in
Determines the AMI id of the desired box version for the configured region.
Writes a file `ami` containing the AMI id.
