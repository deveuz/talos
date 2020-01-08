// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package docker

import (
	"context"
	"net"

	"github.com/talos-systems/talos/internal/pkg/provision"
	"github.com/talos-systems/talos/pkg/config/types/v1alpha1/generate"
)

func (p *provisioner) Reflect(ctx context.Context, clusterName string) (provision.Cluster, error) {
	res := &result{
		clusterInfo: provision.ClusterInfo{
			ClusterName: clusterName,
		},
	}

	// find network assuming network name == cluster name
	networks, err := p.listNetworks(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	if len(networks) > 0 {
		network := networks[0]

		var cidr *net.IPNet
		_, cidr, err = net.ParseCIDR(network.IPAM.Config[0].Subnet)

		if err != nil {
			return nil, err
		}

		res.clusterInfo.Network.Name = network.Name
		res.clusterInfo.Network.CIDR = *cidr
	}

	// find nodes (containers)
	nodes, err := p.listNodes(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		t, err := generate.ParseType(node.Labels["talos.type"])
		if err != nil {
			return nil, err
		}

		res.clusterInfo.Nodes = append(res.clusterInfo.Nodes,
			provision.NodeInfo{
				ID:   node.ID,
				Name: node.Names[0],
				Type: t,

				PrivateIP: net.ParseIP(node.NetworkSettings.Networks[res.clusterInfo.Network.Name].IPAddress),
			})
	}

	return res, nil
}