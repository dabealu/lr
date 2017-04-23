### lr (list registry)  
Simple CLI tool to list and delete images and their tags from docker registry(v2).  
Deletion must be allowed in registry, it can be enabled by running registry container with `REGISTRY_STORAGE_DELETE_ENABLED=true` variable, or by editing registry config.yml (https://docs.docker.com/registry/configuration/#delete).  

**How to use:**  
Set registry connection information inside config file:
```
cat <<EOF > ~/.lr.json
{"addr":"https://registry.example.com","user":"myuser","password":"mypassword"}
EOF
```

or with env variables:
```
export REGISTRY_ADDRESS=https://registry.example.com
export REGISTRY_USER=myuser
export REGISTRY_PASSWORD=mypassword
```
Note that config file takes precedence over env vars.  
  
CLI flags to operate with registry:  
```
ls-images                 # list registry images (short: li)
ls-tags image             # list tags of an image (short: lt)
rm-image image            # remove all image tags (short: ri)
rm-tags image:tag1,tag2   # remove some image tag(s) (short: rt)
help                      # print help
```

**Build binary inside container:**  
```
docker run --rm -v "$PWD":/usr/src/lr -w /usr/src/lr golang:1.8-stretch go build -v
sudo mv lr /usr/local/bin/
```

or build docker image that contains binary:
```
docker build -t lr .
```

and then run lr from container, i.e. list registry images:
```
docker run --rm -e REGISTRY_ADDRESS=https://registry.example.com -e REGISTRY_USER=myuser -e REGISTRY_PASSWORD=mypassword lr li
```

**NOTE:** to free disk space after image/tag deletion, we need to perform garbage collection inside registry container:  
```
registry garbage-collect /etc/docker/registry/config.yml
```
