package allocator

import (
	"context"
	"fmt"
	"net"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusteripv1 "github.com/aojea/clusterip-webhook/api/v1"
	"github.com/go-logr/logr"
)

// copied from k8s.io/kubernetes/pkg/registry/core/service/ipallocator/allocator.go
// Interface manages the allocation of IP addresses out of a range. Interface
// should be threadsafe.
type Interface interface {
	Allocate(net.IP) error
	AllocateNext() (net.IP, error)
	Release(net.IP) error
	ForEach(func(net.IP))
	CIDR() net.IPNet

	// For testing
	Has(ip net.IP) bool
}

// assume that the IPRange object is well-known for the moment
// namespace kube-system name ipv4range

type Range struct {
	client client.Client
	Log    logr.Logger
}

var _ Interface = &Range{}

// NewAllocatorCDRRange creates a Range over a net.IPNet
func NewAllocatorCIDRRange(cidr *net.IPNet, client client.Client) (*Range, error) {
	ctx := context.Background()
	// create IPRange object
	ipRange := clusteripv1.IPRange{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "kube-system",
			Name:      "allocator",
		},
		Spec: clusteripv1.IPRangeSpec{
			Range: cidr.String(),
		},
	}
	err := client.Create(ctx, &ipRange)
	return &Range{
		client: client,
	}, err
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
	addresses := sets.NewString(ipRange.Spec.Addresses...)
	// return if the address doesn't exist in the allocator
	if !addresses.Has(ip.String()) {
		return nil
	}
	addresses.Delete(ip.String())
	ipRange.Spec.Addresses = addresses.List()
	if err := r.client.Update(ctx, ipRange); err != nil {
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
	ctx := context.Background()
	log := r.Log.WithName("iprange")
	ipRange := &clusteripv1.IPRange{}
	key := client.ObjectKey{Namespace: "kube-system", Name: "allocator"}
	if err := r.client.Get(ctx, key, ipRange); err != nil {
		log.Error(err, "unable to fetch IPRange")
		return false
	}
	addresses := sets.NewString(ipRange.Spec.Addresses...)
	return addresses.Has(ip.String())
}
