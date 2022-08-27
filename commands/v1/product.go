package v1

import (
	"fmt"
  "github.com/urfave/cli/v2"
  products "github.com/AlexGustafsson/systembolaget-api/systembolaget/v1/product"
)

// ProductFetchCommand ...
func ProductFetchCommand(context *cli.Context) error {
	err, products := products.Fetch(context.String("key"), context.String("id"))
  return nil
}
