package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/container/v1"
)

func TestAccContainerNodePool_basic(t *testing.T) {
	name := "tf-nodepool-test-" + acctest.RandString(10)
	zone := "us-central1-a"
	clusterConfig := SomeGoogleContainerCluster()
	nodepoolConfig := SomeGoogleContainerNodePool(clusterConfig).
		WithAttribute("name", name).
		WithAttribute("zone", zone).
		WithAttribute("initial_node_count", 2).
		WithAttribute("node_config", NewNestedConfig().
			WithAttribute("machine_type", "n1-highmem-4").
			WithAttribute("labels", NewNestedConfig().
				WithAttribute("my_label", "my_label_value").
				WithAttribute("my_other_label", "my_other_label_value")).
			WithAttribute("tags", "my_tag", "my_other_tag"))

	var nodePool container.NodePool

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: clusterConfig.String() + nodepoolConfig.String(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerNodePoolExists(zone, clusterConfig.Name(), name, &nodePool),
					testAccCheckContainerNodePoolHasInitialNodeCount(&nodePool, 2),
					testAccCheckContainerNodePoolHasLabel(&nodePool, "my_label", "my_label_value"),
					testAccCheckContainerNodePoolHasLabel(&nodePool, "my_other_label", "my_other_label_value"),
					testAccCheckContainerNodePoolHasTags(&nodePool, "my_tag", "my_other_tag"),
					testAccCheckContainerNodePoolHasMachineType(&nodePool, "n1-highmem-4"),
				),
			},
		},
	})
}

func testAccCheckContainerNodePoolDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_container_node_pool" {
			continue
		}

		attributes := rs.Primary.Attributes
		_, err := config.clientContainer.Projects.Zones.Clusters.NodePools.Get(
			config.Project, attributes["zone"], attributes["cluster"], attributes["name"]).Do()
		if err == nil {
			return fmt.Errorf("NodePool still exists")
		}
	}

	return nil
}

func SomeGoogleContainerCluster() *ConfigBuilder {
	return NewResourceConfigBuilder("google_container_cluster", "cluster-"+acctest.RandString(10)).
		WithAttribute("name", "tf-cluster-nodepool-test-"+acctest.RandString(10)).
		WithAttribute("zone", "us-central1-a").
		WithAttribute("initial_node_count", 3).
		WithAttribute("master_auth", NewNestedConfig().
			WithAttribute("username", "mr.yoda").
			WithAttribute("password", "adoy.rm"))
}

func SomeGoogleContainerNodePool(cluster *ConfigBuilder) *ConfigBuilder {
	return NewResourceConfigBuilder("google_container_node_pool", "nodepool-"+acctest.RandString(10)).
		WithAttribute("name", "tf-nodepool-test-"+acctest.RandString(10)).
		WithAttribute("zone", "us-central1-a").
		WithAttribute("cluster", fmt.Sprintf("${google_container_cluster.%s.name}", cluster.ResourceName)).
		WithAttribute("initial_node_count", 2)
}

func testAccCheckContainerNodePoolExists(zone, clusterName, nodePoolName string, nodePool *container.NodePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		found, err := config.clientContainer.Projects.Zones.Clusters.NodePools.Get(config.Project, zone, clusterName, nodePoolName).Do()
		if err != nil {
			return err
		}

		if found == nil {
			return fmt.Errorf("Unable to find resource")
		}

		*nodePool = *found
		return nil
	}
}

func testAccCheckContainerNodePoolHasInitialNodeCount(nodePool *container.NodePool, count int64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nodePool.InitialNodeCount != count {
			return fmt.Errorf("Expected initial_node_count %d but found %d", count, nodePool.InitialNodeCount)
		}
		return nil
	}
}

func testAccCheckContainerNodePoolHasLabel(nodePool *container.NodePool, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := nodePool.Config.Labels[key]
		if !ok {
			return fmt.Errorf("Label %s not found", key)
		}

		if v != value {
			return fmt.Errorf("Label key '%s' has incorrect value; expected '%s', found '%s'", key, value, v)
		}
		return nil
	}
}

func testAccCheckContainerNodePoolHasMachineType(nodePool *container.NodePool, machineType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nodePool.Config.MachineType != machineType {
			return fmt.Errorf("Expected machine type '%s', found '%s'", machineType, nodePool.Config.MachineType)
		}

		return nil
	}
}

func testAccCheckContainerNodePoolHasTags(nodePool *container.NodePool, tags ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Perform simple set difference
		tagsNeeded := map[string]interface{}{}
		for _, x := range tags {
			tagsNeeded[x] = struct{}{}
		}

		for _, x := range nodePool.Config.Tags {
			delete(tagsNeeded, x)
		}

		if len(tagsNeeded) > 0 {
			// Convert to array for pretty print
			missing := make([]string, len(tagsNeeded))
			idx := 0
			for key := range tagsNeeded {
				missing[idx] = key
				idx++
			}

			return fmt.Errorf("Tags not found: %v", missing)
		}

		return nil
	}
}
