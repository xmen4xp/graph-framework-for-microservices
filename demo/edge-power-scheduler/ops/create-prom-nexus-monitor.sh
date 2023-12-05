NAMESPACE="${NAMESPACE:=default}"
CACRT=`kubectl -n ${NAMESPACE} get secret  nexus-k8s-apiserver-kubeconfig -o json | jq -r '.data.kubeconfig' | base64 -d | yq '.users.[].user.client-certificate-data'`
CAKEY=`kubectl -n ${NAMESPACE} get secret  nexus-k8s-apiserver-kubeconfig -o json | jq -r '.data.kubeconfig' | base64 -d | yq '.users.[].user.client-key-data'`

kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: nexus-prom-k8s-api-cert
  namespace: prometheus
data:
  ca.crt: ${CACRT}
  ca.key: ${CAKEY}
EOF

kubectl apply -f - <<EOF
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  labels:
    app: kube-prometheus-stack-nexus-k8s-apiserver
    release: stable
  name: stable-kube-prometheus-sta-nexus-k8s-apiserver
  namespace: prometheus
spec:
  endpoints:
  - bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
    metricRelabelings:
    - action: drop
      regex: apiserver_request_duration_seconds_bucket;(0.15|0.2|0.3|0.35|0.4|0.45|0.6|0.7|0.8|0.9|1.25|1.5|1.75|2|3|3.5|4|4.5|6|7|8|9|15|25|40|50)
      sourceLabels:
      - __name__
      - le
    port: https
    scheme: https
    tlsConfig:
      cert: 
        secret:
          name: nexus-prom-k8s-api-cert
          key: ca.crt
      keySecret:
        name: nexus-prom-k8s-api-cert
        key: ca.key        
      insecureSkipVerify: true
      serverName: kubernetes
  jobLabel: component
  namespaceSelector:
    matchNames:
    - ${NAMESPACE}
  selector:
    matchLabels:
      component: nexus-apiserver
EOF

kubectl apply -f - <<EOF
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: stable-kube-prometheus-sta-nexus-etcd
  namespace: prometheus
  labels:
    app: kube-prometheus-sta-nexus-etcd
    release: stable
spec:
  endpoints:
  - bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
    port: client
  jobLabel: jobLabel
  namespaceSelector:
    matchNames:
    - ${NAMESPACE}
  selector:
    matchLabels:
      app: nexus-etcd
EOF