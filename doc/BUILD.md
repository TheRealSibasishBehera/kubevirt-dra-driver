# How to build Kubevirt DRA Driver container image

## Platforms supported

- Linux

## Prerequisites

- Docker

## Building

- Script to simply rebuild the image is already present in demo folder
 ```bash
cd demo
./build-driver.sh
```
- If you want to rebuild the CRD specifically to further change logic in the driver you can use the `Makefile`

`Makefile` automates this, only required tool is Docker , where the generation process takes place in a container and is copied to the configured path .

```bash
make docker-generate
```
- If you want to rebuild the whole build

```bash
make docker-build
```

