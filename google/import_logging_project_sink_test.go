package google

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccLoggingProjectSink_importBasic(t *testing.T) {
	sinkName := "tf-test-" + acctest.RandString(10)
	bucketName := "tf-test-" + acctest.RandString(10)

	resourceName := "google_logging_project_sink.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingProjectSinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_basic(sinkName, bucketName),
			}, {
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
