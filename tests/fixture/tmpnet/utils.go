// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package tmpnet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	DefaultNodeTickerInterval = 50 * time.Millisecond
)

var ErrNotRunning = errors.New("not running")

// WaitForHealthy blocks until Node.IsHealthy returns true or an error (including context timeout) is observed.
func WaitForHealthy(ctx context.Context, node *Node) error {
	if _, ok := ctx.Deadline(); !ok {
		return fmt.Errorf("unable to wait for health for node %q with a context without a deadline", node.NodeID)
	}
	ticker := time.NewTicker(DefaultNodeTickerInterval)
	defer ticker.Stop()

	for {
		healthy, err := node.IsHealthy(ctx)
		if err != nil && !errors.Is(err, ErrNotRunning) {
			return fmt.Errorf("failed to wait for health of node %q: %w", node.NodeID, err)
		}
		if healthy {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("failed to wait for health of node %q before timeout: %w", node.NodeID, ctx.Err())
		case <-ticker.C:
		}
	}
}

// Marshal to json with default prefix and indent.
func DefaultJSONMarshal(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}