package end2end

import (
	"context"

	"github.com/Velocidex/ordereddict"
)

// Add checks that the server is good.
func CheckServer(
	ctx context.Context, state *ordereddict.Dict,
	binary, api_config_file string) error {
	return CheckLabels(ctx, state, binary, api_config_file)
}
