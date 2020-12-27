/*
Copyright 2020.

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
	"context"
	"net"

	"github.com/go-logr/logr"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusteripv1 "github.com/aojea/clusterip-webhook/api/v1"
)

// ServiceReconciler reconciles a Service object
type ServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=clusterip.allocator.x-k8s.io,resources=ipranges,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=clusterip.allocator.x-k8s.io,resources=ipranges/status,verbs=get;update;patch

func (r *ServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("service", req.NamespacedName)

	// get all services
	var svcList v1.ServiceList
	if err := r.List(ctx, &svcList); err != nil {
		log.Error(err, "unable to list services")
		return ctrl.Result{}, err
	}
	// obtain all assigned clusterIPs
	svcIPs := sets.NewString()
	for _, svc := range svcList.Items {
		ip := net.ParseIP(svc.Spec.ClusterIP)
		if ip != nil {
			svcIPs.Insert(svc.Spec.ClusterIP)
		}

	}
	// obtain current allocator addresses
	ipRange := &clusteripv1.IPRange{}
	key := client.ObjectKey{Namespace: "kube-system", Name: "allocator"}
	if err := r.Get(ctx, key, ipRange); err != nil {
		log.Error(err, "unable to fetch IPRange")
		return ctrl.Result{}, err
	}
	addresses := sets.NewString(ipRange.Spec.Addresses...)
	// reconcile
	diff := svcIPs.Difference(addresses)
	if len(diff) > 0 {
		log.Info("allocator is not synced", diff)
	}
	ipRange.Spec.Addresses = svcIPs.List()
	if err := r.Update(ctx, ipRange); err != nil {
		log.Error(err, "unable to update IPRange")
		return ctrl.Result{}, err
	}
	//
	return ctrl.Result{}, nil
}

func (r *ServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Service{}).
		Complete(r)
}
