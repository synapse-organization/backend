cd ..
go mod vendor
#docker rmi backend:latest
docker build . -t irania9o/synapse_backend -f deploy/Dockerfile
docker push irania9o/synapse_backend
#docker save -o deploy/backend.tar backend:latest