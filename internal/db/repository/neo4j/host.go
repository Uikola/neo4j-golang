package neo4j

import (
	"context"
	"fmt"
	"github.com/Uikola/neo4j-golang/internal/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Repository struct {
	driver neo4j.DriverWithContext
}

func NewRepository(driver neo4j.DriverWithContext) *Repository {
	return &Repository{driver: driver}
}

func (r Repository) Save(ctx context.Context, host entity.Host) error {
	const op = "Neo4j.HostRepository.Save"

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	interfaces := make([]map[string]interface{}, len(host.Interfaces))
	for i, iface := range host.Interfaces {
		interfaces[i] = map[string]interface{}{
			"name":   iface.Name,
			"mtu":    iface.MTU,
			"rx_pkt": iface.RxPackets,
			"rx_bts": iface.RxBytes,
			"tx_pkt": iface.TxPackets,
			"tx_bts": iface.TxBytes,
		}
	}

	cypherQuery := `
	MERGE (h:Host {hostname: $hostname})
	WITH h
	UNWIND $interfaces AS iface
	MERGE (i:Interface {name: iface.name})
	SET i.mtu = iface.mtu,
	    i.rx_pkt = iface.rx_pkt,
	    i.rx_bts = iface.rx_bts,
	    i.tx_pkt = iface.tx_pkt,
	    i.tx_bts = iface.tx_bts
	MERGE (h)-[:HAS_NETWORK_INTERFACE]->(i)
	`

	params := map[string]interface{}{
		"hostname":   host.Hostname,
		"interfaces": interfaces,
	}

	_, err := session.Run(ctx, cypherQuery, params)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	return nil
}
