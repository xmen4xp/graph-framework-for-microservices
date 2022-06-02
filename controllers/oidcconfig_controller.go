/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"api-gw/pkg/model"
	"context"
	log "github.com/sirupsen/logrus"
	authnexusv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/api.git/build/apis/authentication.nexus.org/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// OidcConfigReconciler reconciles a OidcConfig object
type OidcConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=authentication.nexus.org.api-gw.com,resources=oidcconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=authentication.nexus.org.api-gw.com,resources=oidcconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=authentication.nexus.org.api-gw.com,resources=oidcconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the OidcConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *OidcConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Business Logic: An event has occurred on the OIDC Node.
	var oidcNode authnexusv1.OIDC
	eventType := model.Upsert
	// TODO use the shim layer to fetch the object
	if err := r.Get(ctx, req.NamespacedName, &oidcNode); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Errorf("Error while trying to fetch OIDC node with name %s\n", req.Name)
			return ctrl.Result{}, err
		}
		eventType = model.Delete
	}
	log.Debugf("Received event %s for oidcNode node: Name %s\n", eventType, oidcNode.Name)

	// Pass on the event to the authenticator so that it can reconfigure itself
	model.OidcChan <- model.OidcNodeEvent{Oidc: oidcNode, Type: eventType}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OidcConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&authnexusv1.OIDC{}).
		Complete(r)
}
