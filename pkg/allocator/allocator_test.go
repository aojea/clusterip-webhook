package allocator

import (
	"context"
	"net"
	"path/filepath"
	"testing"

	clusteripv1 "github.com/aojea/clusterip-webhook/api/v1"

	"k8s.io/client-go/deprecated/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func TestMain(t *testing.T) {
	key := client.ObjectKey{Namespace: "kube-system", Name: "allocator"}
	ctx := context.Background()
	ipRange := &clusteripv1.IPRange{}
	// specify testEnv configuration
	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("../..", "config", "crd", "bases")},
	}

	// start testEnv
	cfg, err := testEnv.Start()
	if err != nil {
		t.Fatalf("Unable to start test environment: (%v)", err)
	}

	// Add clusterIP allocator to scheme
	if err := clusteripv1.AddToScheme(scheme.Scheme); err != nil {
		t.Fatalf("Unable to add iprange scheme: (%v)", err)
	}

	// +kubebuilder:scaffold:scheme

	cs, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		t.Fatalf(err.Error())
	}
	ip, subnet, _ := net.ParseCIDR("10.96.0.2/24")
	r, err := NewAllocatorCIDRRange(subnet, cs)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if err := cs.Get(ctx, key, ipRange); err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("new iprange object %v", ipRange)
	// Allocate an specific IP address
	err = r.Allocate(ip)
	if err != nil {
		t.Fatalf(err.Error())
	}
	for i := 0; i < 254; i++ {
		alloc, err := r.AllocateNext()
		if err != nil {
			t.Fatalf(err.Error())
		}
		t.Logf("allocated %v", alloc)
	}
	// Check the new object
	if err := cs.Get(ctx, key, ipRange); err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("new Range addresses %v", ipRange.Spec.Addresses)

	// stop testEnv
	err = testEnv.Stop()
}
