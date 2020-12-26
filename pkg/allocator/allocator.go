package allocator

import (
	"context"
	"fmt"
	"net"

	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusteripv1 "github.com/aojea/clusterip-webhook/api/v1"
	"github.com/go-logr/logr"
)

// assume that the IPRange object is well-known for the moment
// namespace kube-system name ipv4range

type Range struct {
	client client.Client
	Log    logr.Logger
}

func (r *Range) Allocate(ip net.IP) error {
	ctx := context.Background()
	log := r.Log.WithValues("iprange", ip)
	ipRange := &clusteripv1.IPRange{}
	key := client.ObjectKey{Namespace: "kube-system", Name: "allocator"}
	if err := r.client.Get(ctx, key, ipRange); err != nil {
		log.Error(err, "unable to fetch IPRange")
		return nil
	}
	addresses := sets.NewString(ipRange.Spec.Addresses...)
	if addresses.Has(ip.String()) {
		return fmt.Errorf("ip %s already allocated", ip.String())
	}
	addresses.Insert(ip.String())
	ipRange.Spec.Addresses = addresses.List()
	if err := r.client.Update(ctx, ipRange); err != nil {
		log.Error(err, "unable to fetch IPRange")
		return err
	}
	return nil
}

func (r *Range) AllocateNext() (net.IP, error) {
	ctx := context.Background()
	log := r.Log.WithName("iprange")
	ipRange := &clusteripv1.IPRange{}
	key := client.ObjectKey{Namespace: "kube-system", Name: "allocator"}
	if err := r.client.Get(ctx, key, ipRange); err != nil {
		log.Error(err, "unable to fetch IPRange")
		return nil, err
	}
	// Range is validated by the webhook
	ip, _, _ := net.ParseCIDR(ipRange.Spec.Range)
	return ip, nil
}

func (r *Range) Release(ip net.IP) error {
	ctx := context.Background()
	log := r.Log.WithValues("iprange", ip)
	ipRange := &clusteripv1.IPRange{}
	key := client.ObjectKey{Namespace: "kube-system", Name: "allocator"}
	if err := r.client.Get(ctx, key, ipRange); err != nil {
		log.Error(err, "unable to fetch IPRange")
		return err
	}
	return nil
}

func (r *Range) ForEach(func(net.IP)) {

}

func (r *Range) CIDR() net.IPNet {
	ctx := context.Background()
	log := r.Log.WithName("iprange")
	ipRange := &clusteripv1.IPRange{}
	key := client.ObjectKey{Namespace: "kube-system", Name: "allocator"}
	if err := r.client.Get(ctx, key, ipRange); err != nil {
		log.Error(err, "unable to fetch IPRange")
		return net.IPNet{}
	}
	// Range is validated by the webhook
	_, subnet, _ := net.ParseCIDR(ipRange.Spec.Range)
	return *subnet

}

// For testing
func (r *Range) Has(ip net.IP) bool {
	return true
}
