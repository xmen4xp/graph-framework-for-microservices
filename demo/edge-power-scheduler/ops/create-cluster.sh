set -e
set -x
export EKS_CLUSTER_NAME=$1
if [[ -z "$EKS_CLUSTER_NAME" ]]; then
    echo "run with cluster name as argument"
    exit 1
fi
export EKS_REGION='us-west-2'
export ROLE_NAME="${EKS_CLUSTER_NAME}_EBS_CSI_DriverRole"
export ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text)
export ACCOUNT_ROLE_ARN="arn:aws:iam::$ACCOUNT_ID:role/$ROLE_NAME"
export POLICY_ARN="arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy"
echo "Creating EKS cluster with name:${EKS_CLUSTER_NAME} in region:${EKS_REGION}"
# User running this command should have iam permission as per https://eksctl.io/usage/minimum-iam-policies/
# on our account a group exists "eks-group" just add user to that group to get the permissions

cat <<EOF > eks-${EKS_CLUSTER_NAME}.yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: ${EKS_CLUSTER_NAME}
  region: ${EKS_REGION}
  tags:
    creator: ${USER}
    environment: dev
iam:
  withOIDC: true
vpc: 
  id: "vpc-0e74325a800180fb3"
  subnets:
    private:
      us-west-2a: {id: subnet-09393ebe114a268c0 }
      us-west-2b: {id: subnet-0a378d1e32a84e7cc }
      us-west-2c: {id: subnet-01036423c0e4e72b3 }
      us-west-2d: {id: subnet-03ef2d408af1feb90 }
managedNodeGroups:
  - name: ${EKS_CLUSTER_NAME}-ng-0
    instanceType: c5.4xlarge
    desiredCapacity: 3
    volumeSize: 80
    labels: { role: workers, nexus: api }
    taints:
      - {
          "key": "nexus",
          "value": "api",
          "effect": "NoSchedule"
        }
    privateNetworking: true
    iam:
      withAddonPolicies:
        imageBuilder: true
        autoScaler: true
        ebs: true
        efs: true
        fsx: true
        cloudWatch: true
        externalDNS: true
        certManager: true
        albIngress: true
    ssh:
      allow: true
  - name: ${EKS_CLUSTER_NAME}-ng-1
    instanceType: c5.2xlarge
    desiredCapacity: 8
    volumeSize: 80
    labels: { role: workers }
    privateNetworking: true
    iam:
      withAddonPolicies:
        imageBuilder: true
        autoScaler: true
        ebs: true
        efs: true
        fsx: true
        cloudWatch: true
        externalDNS: true
        certManager: true
        albIngress: true
    ssh:
      allow: true
  - name: ${EKS_CLUSTER_NAME}-ng-2
    instanceType: c5.2xlarge
    desiredCapacity: 3
    volumeSize: 80
    labels: { role: workers, eps: scheduler }
    taints:
      - {
          "key": "eps",
          "value": "scheduler",
          "effect": "NoSchedule"
        }
    privateNetworking: true
    iam:
      withAddonPolicies:
        imageBuilder: true
        autoScaler: true
        ebs: true
        efs: true
        fsx: true
        cloudWatch: true
        externalDNS: true
        certManager: true
        albIngress: true
    ssh:
      allow: true
  - name: ${EKS_CLUSTER_NAME}-ng-3
    instanceType: c5.2xlarge
    desiredCapacity: 3
    volumeSize: 80
    labels: { role: workers, eps: generator }
    taints:
      - {
          "key": "eps",
          "value": "generator",
          "effect": "NoSchedule"
        }
    privateNetworking: true
    iam:
      withAddonPolicies:
        imageBuilder: true
        autoScaler: true
        ebs: true
        efs: true
        fsx: true
        cloudWatch: true
        externalDNS: true
        certManager: true
        albIngress: true
    ssh:
      allow: true
EOF

eksctl create cluster -f ./eks-${EKS_CLUSTER_NAME}.yaml

export oidc_id=$(aws eks describe-cluster --name $EKS_CLUSTER_NAME --query "cluster.identity.oidc.issuer" --output text | cut -d '/' -f 5)


# AWS IAM role bound to a Kubernetes service account
eksctl create iamserviceaccount \
  --name "ebs-csi-controller-sa" \
  --namespace "kube-system" \
  --cluster $EKS_CLUSTER_NAME \
  --region $EKS_REGION \
  --attach-policy-arn $POLICY_ARN \
  --role-only \
  --role-name $ROLE_NAME \
  --approve

# Create Addon
eksctl create addon \
  --name "aws-ebs-csi-driver" \
  --cluster "$EKS_CLUSTER_NAME" \
  --region="$EKS_REGION" \
  --service-account-role-arn $ACCOUNT_ROLE_ARN \
  --force

# Get status of the driver, must be STATUS=ACTIVE
eksctl get addon \
  --name "aws-ebs-csi-driver" \
  --region $EKS_REGION \
  --cluster "$EKS_CLUSTER_NAME"

