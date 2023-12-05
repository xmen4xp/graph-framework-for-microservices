

# Installing prometheus and grafana

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install stable prometheus-community/kube-prometheus-stack -n prometheus
# To have the prometheus monitor your nexus etcd and api server
./create-prom-nexus-monitor.sh


## log into grafana
kubectl -n prometheus port-forward svc/stable-grafana 8800:80
### open browser page at localhost:8800
### username is admin default p is prom-operator
### you can get default values by running 
helm show values prometheus-community/kube-prometheus-stack > kube-prometheus-stack.values
# to change replica count
k patch statefulsets eps-edge-agent -p '{"spec":{"replicas":10}}'


