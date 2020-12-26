package allocator

import (
	"math/big"
	"net"

	"k8s.io/kubernetes/pkg/registry/core/service/ipallocator"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusteripv1 "github.com/aojea/clusterip-webhook/api/v1"
)

// assume that the IPRange object is well-known for the moment
// namespace kube-system name ipv4range

type Range struct {
	client client.Client
	Log    logr.Logger
}

var _ ipallocator.Interface = &Range{}

func (r *Range) Allocate(ip net.IP) error {
	ctx := context.Background()
	log := r.Log.WithValues("iprange", ip)
	var ipRange &clusteripv1.IPRange
	ipRange.Spec.Addresses = ip.String()
	r.Patch(ctx, )
}

func (r *Range) AllocateNext() (net.IP, error) {
	ctx := context.Background()
	var ipRange &clusteripv1.IPRange
    if err := r.Get(ctx, req.NamespacedName, &ipRange); err != nil {
        log.Error(err, "unable to fetch IPRange")
        return nil
	}
	// Range is validated by the webhook
	_, subnet, _ := net.ParseCIDR(ipRange.Spec.Range)
	return subnet
}

func (r *Range) Release(net.IP) error {
	ctx := context.Background()
	log := r.Log.WithValues("iprange", ip)
	var ipRange &clusteripv1.IPRange
	ipRange.Spec.Addresses = ip.String()
	r.Patch(ctx, )
}

func (r *Range) ForEach(func(net.IP)) {

}

func (r *Range) CIDR() net.IPNet {
	ctx := context.Background()
	var ipRange &clusteripv1.IPRange
    if err := r.Get(ctx, req.NamespacedName, &ipRange); err != nil {
        log.Error(err, "unable to fetch IPRange")
        return nil
	}
	// Range is validated by the webhook
	_, subnet, _ := net.ParseCIDR(ipRange.Spec.Range)
	return subnet

}

// For testing
func (r *Range) Has(ip net.IP) bool{

}