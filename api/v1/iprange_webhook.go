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

package v1

import (
	"fmt"
	"net"

	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var iprangelog = logf.Log.WithName("iprange-resource")

func (r *IPRange) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-clusterip-allocator-x-k8s-io-v1-iprange,mutating=true,failurePolicy=fail,groups=clusterip.allocator.x-k8s.io,resources=ipranges,verbs=create;update,versions=v1,name=miprange.kb.io

var _ webhook.Defaulter = &IPRange{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *IPRange) Default() {
	iprangelog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// +kubebuilder:webhook:verbs=create;update;delete,path=/validate-clusterip-allocator-x-k8s-io-v1-iprange,mutating=false,failurePolicy=fail,groups=clusterip.allocator.x-k8s.io,resources=ipranges,versions=v1,name=viprange.kb.io

var _ webhook.Validator = &IPRange{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *IPRange) ValidateCreate() error {
	iprangelog.Info("validate create", "name", r.Name)
	// Create only allows to set the IP range
	if len(r.Spec.Addresses) > 0 {
		return fmt.Errorf("Addresses can not be allocated on creation")
	}
	_, ipRange, err := net.ParseCIDR(r.Spec.Range)
	if err != nil {
		return err
	}
	r.Spec.Range = ipRange.String()
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *IPRange) ValidateUpdate(old runtime.Object) error {
	oldIPRange := old.(*IPRange)
	iprangelog.Info("validate update", "name", r.Name)
	// Range is inmutable after creation
	if r.Spec.Range != oldIPRange.Spec.Range {
		return fmt.Errorf("Range can not be changed after creation")
	}
	// the Range was already validated on creation
	_, ipRange, _ := net.ParseCIDR(r.Spec.Range)
	allErrors := []error{}
	for _, address := range r.Spec.Addresses {
		ip := net.ParseIP(address)
		if ip == nil {
			allErrors = append(allErrors, fmt.Errorf("invalid ip address %s", address))
		}
		if !ipRange.Contains(ip) {
			allErrors = append(allErrors, fmt.Errorf("ip address %s out of range %s", address, ipRange.String()))
		}
		if ip.Equal(ipRange.IP) {
			allErrors = append(allErrors, fmt.Errorf("ip address %s reserved", ip.String()))
		}
	}
	return utilerrors.NewAggregate(allErrors)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *IPRange) ValidateDelete() error {
	iprangelog.Info("validate delete", "name", r.Name)
	// An IPRange can not be deleted if there are still ip addresses allocated
	if len(r.Spec.Addresses) > 0 {
		return fmt.Errorf("IPRange can not be deleted if addresses are allocated")
	}
	return nil
}
