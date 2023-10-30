#!/bin/bash

set -x

cluster_name=""
delete_cluster=false
start_port=""

# cleanup_cluster handles deletion of all resources related to kind setup
function cleanup_cluster {
    rm -rf $CLUSTER_DIR
    kind delete cluster -n $cluster_name
    docker rm -f k8s-proxy-$cluster_name
}

while getopts "n:ds:m:h:k:p:" flag; do
    case "${flag}" in
        n) cluster_name=${OPTARG};;
        d) delete_cluster=true;;
        p) start_port=${OPTARG};;
        s) server_port=${OPTARG};;
        m) metrics_address=${OPTARG};;
        h) health_address=${OPTARG};;
        k) kubctl_proxy_port=${OPTARG};;
    esac
done

if [ -z "$cluster_name" ]
then
      echo "cluster name is mandatory. Provide a valid cluster name with -n argument"
      exit 1
fi

if [ "$delete_cluster" == true ] ; then
    cleanup_cluster
    exit 0
fi

if [ -z "$start_port" ]
then
      echo "start port number is mandatory. Provide a valid start port number with -s argument"
      exit 1
fi

# kubectl proxy port
kubctl_proxy_port=$start_port

next_port=$((start_port + 1))
if [ -z "$server_port" ]
then
      server_port=$next_port
      next_port=$((next_port + 1))
fi

if [ -z "$metrics_address" ]
then
      metrics_address=":$next_port"
      next_port=$((next_port + 1))
fi

if [ -z "$health_address" ]
then
      health_address=":$next_port"
      next_port=$((next_port + 1))
fi

if [ -z "$cluster_name" ]
then
      echo "cluster name is mandatory. Provide a valid cluster name with -n argument"
      exit 1
fi

# Setup up directory related variables.
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
CLUSTER_DIR=${SCRIPT_DIR}/"."${cluster_name}

if [ -d "$CLUSTER_DIR" ]; then
    cleanup_cluster
fi

# create cluster directory
mkdir -p $CLUSTER_DIR

sed -e "s/HTTP_SERVER_PORT/$server_port/g" -e "s/HEALTH_PROBE_ADDRESS/$health_address/g" -e "s/METRICS_ADDRESS/$metrics_address/g" $SCRIPT_DIR/api-gw-config.tmpl > $CLUSTER_DIR/api-gw-config.yaml

# Create the kind k8s cluster
KUBECONFIG_FILE=${CLUSTER_DIR}/kubeconfig
kind create cluster -n $cluster_name
kind export kubeconfig -n $cluster_name --kubeconfig ${KUBECONFIG_FILE}
chmod 666 ${KUBECONFIG_FILE}

MOUNTED_KUBECONFIG=/etc/config/kubeconfig
docker run -d --name=k8s-proxy-$cluster_name --rm --network host --pull=missing --mount type=bind,source=${KUBECONFIG_FILE},target=${MOUNTED_KUBECONFIG},readonly -e KUBECONFIG=${MOUNTED_KUBECONFIG} bitnami/kubectl proxy -p $kubctl_proxy_port --disable-filter=true --v=1

exit 0
