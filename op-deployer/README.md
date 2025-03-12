
# Custom image with predeployed contracts support

```bash
# Build static binary for the docker image
go build -o op-deployer -ldflags="-extldflags=-static" cmd/op-deployer/main.go

# Dockerfile for the new image just overwriting binary
>Dockerfile.tmp cat <<EOF
FROM us-docker.pkg.dev/oplabs-tools-artifacts/images/op-deployer:v0.0.12
COPY op-deployer /usr/local/bin/op-deployer
CMD ["op-deployer"]
EOF

# Build the new image, tag and push as required
docker build -t op-deployer . -f Dockerfile.tmp && rm Dockerfile.tmp
docker tag op-deployer:latest whatever/op-deployer:whatever
docker push whatever/op-deployer:whatever

```
