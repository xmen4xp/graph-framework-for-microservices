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
echo "Delating EKS cluster with name:${EKS_CLUSTER_NAME} in region:${EKS_REGION}"
echo "Deleting clusters will not clean up volums that where attached. please chean form console."
# AWS IAM role bound to a Kubernetes service account
eksctl delete iamserviceaccount \
  --name "ebs-csi-controller-sa" \
  --namespace "kube-system" \
  --cluster $EKS_CLUSTER_NAME \
  --region $EKS_REGION 

eksctl delete cluster --name ${EKS_CLUSTER_NAME}


