Datamodel build:

make runtime.install
export HOST_KUBECONFIG
export HOST_KUBECONFIG_LOCAL

make datamodel_build
make dm.install

make runtime.install.docker

watch -n3 "curl -s -m 60 'localhost:8080/systemstatus'"

curl -s localhost:8080/systemstatus
curl -s -m 60 localhost:8080/broadcast
curl -s localhost:8080/lightcoldboot


curl light-connector-0.light-connector:8080/clear 
curl localhost:8080/broadcast
curl localhost:8080/broadcastinfo
curl -m 60 localhost:8080/broadcast

To run tests:
-------------

cd traffic-light-test
go test ./...
