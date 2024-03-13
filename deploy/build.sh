cd ..
go mod vendor
docker rmi backend:latest
docker build . -t backend:latest -f deploy/Dockerfile
docker save -o deploy/backend.tar backend:latest