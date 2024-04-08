cd ..
go mod vendor
#docker rmi backend:latest
docker build . -t irania9o/synapse_backend:latest -f deploy/Dockerfile
docker push irania9o/synapse_backend:latest
#docker save -o deploy/backend.tar backend:latest