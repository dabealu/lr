### lr (list registry)  
Simple CLI tool to list and delete images and their tags from docker registry(v2).  
Deletion must be allowed in registry, it can be enabled by running registry container with `REGISTRY_STORAGE_DELETE_ENABLED=true` variable, or by editing registry config.yml (https://docs.docker.com/registry/configuration/#delete).  

**How to use:**  
Set registry connection information:  
```
export REGISTRY_ADDRESS=https://registry.example.com
export REGISTRY_USER=user
export REGISTRY_PASSWORD=password
```
CLI flags to operate with registry:  
```
ls-images                 # list registry images (short: li)
ls-tags image             # list tags of an image (short: lt)
rm-image image            # remove all image tags (short: ri)
rm-tags image:tag1,tag2   # remove some image tag(s) (short: rt)
help                      # print help
```
**Build inside container:**  
```
docker run --rm -v "$PWD":/usr/src/lr -w /usr/src/lr golang:1.8-stretch go build -v
sudo mv lr /usr/local/bin/
```
or build docker image that contains binary:
```
docker build -t lr .
```
Run lr from container, i.e. list registry images:
```
docker run --rm -e REGISTRY_ADDRESS=https://registry.example.com -e REGISTRY_USER=usr -e REGISTRY_PASSWORD=passwd lr li
```
