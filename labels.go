package end2end

import (
	"context"
	"errors"
	"fmt"

	"github.com/Velocidex/ordereddict"
)

// Add some labels to the clients
func AddLabels(
	ctx context.Context, state *ordereddict.Dict, binary, api_config_file string) error {
	rows, err := Query(ctx, binary, api_config_file,
		"SELECT client_id, label(client_id=client_id, op='set', labels='Label1') FROM clients() WHERE os_info.hostname =~ '-1$'")
	if err != nil {
		return err
	}
	if len(rows) != 1 {
		return errors.New("Unexpected number of clients updated!")
	}

	client_id, _ := rows[0].GetString("client_id")
	state.Set("Label1", client_id)

	rows, err = Query(ctx, binary, api_config_file,
		"SELECT client_id, label(client_id=client_id, op='set', labels='Label2') FROM clients() WHERE os_info.hostname =~ '-2$'")
	if err != nil {
		return err
	}
	if len(rows) != 1 {
		return errors.New("Unexpected number of clients updated!")
	}
	client_id, _ = rows[0].GetString("client_id")
	state.Set("Label2", client_id)

	return nil
}

func CheckLabels(ctx context.Context, state *ordereddict.Dict,
	binary, api_config_file string) error {

	// Make sure searching the client label index works.
	rows, err := Query(ctx, binary, api_config_file,
		"SELECT client_id FROM clients(search='label:Label1')")
	if err != nil {
		return err
	}
	if len(rows) != 1 {
		return fmt.Errorf("Unexpected label search %v", rows)
	}
	return nil
}
