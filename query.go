package end2end

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Velocidex/ordereddict"
)

func Query(ctx context.Context, binary, api_config_file, query_str string) (
	[]*ordereddict.Dict, error) {
	return query(ctx, binary, api_config_file, query_str, false)
}

func QueryVerbose(ctx context.Context, binary, api_config_file, query_str string) (
	[]*ordereddict.Dict, error) {
	return query(ctx, binary, api_config_file, query_str, true)
}

func query(ctx context.Context,
	binary, api_config_file, query_str string,
	verbose bool) ([]*ordereddict.Dict, error) {

	argv := []string{"--api_config", api_config_file, "query", query_str}
	if verbose {
		argv = append(argv, "-v")
	}

	fmt.Printf("Running %v %v\n", binary, strings.Join(argv, " "))

	cmd := exec.CommandContext(ctx, binary, argv...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error %v\n", string(out))
		return nil, err
	}

	rows := []json.RawMessage{}
	err = json.Unmarshal(out, &rows)
	if err != nil {
		fmt.Printf("Query %v: %v\n", err, string(out))
		return nil, err
	}

	result := make([]*ordereddict.Dict, 0, len(rows))
	for _, row := range rows {
		item := ordereddict.NewDict()
		err := item.UnmarshalJSON(row)
		if err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}
