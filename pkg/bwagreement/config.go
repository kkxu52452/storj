// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package bwagreement

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	monkit "gopkg.in/spacemonkeygo/monkit.v2"
	"storj.io/storj/pkg/pb"
	"storj.io/storj/pkg/provider"
)

var (
	mon = monkit.Package()
)

// Config is a configuration struct that is everything you need to start an
// agreement receiver responsibility
type Config struct {
	DatabaseURL    string `help:"the database connection string to use" default:"postgres://postgres@localhost/pointerdb?sslmode=disable"`
	DatabaseDriver string `help:"the database driver to use" default:"postgres"`
}

// Run implements the provider.Responsibility interface
func (c Config) Run(ctx context.Context, server *provider.Provider) (err error) {
	defer mon.Task()(&ctx)(&err)
	k := server.Identity().Leaf.PublicKey

	ns, err := NewServer(ctx, c.DatabaseDriver, c.DatabaseURL, zap.L(), k)
	if err != nil {
		return fmt.Errorf("Can't connect to the database for bwagreement: %s", err)
	}

	pb.RegisterBandwidthServer(server.GRPC(), ns)

	return server.Run(ctx)
}
