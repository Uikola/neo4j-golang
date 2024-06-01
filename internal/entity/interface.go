package entity

type Interface struct {
	Name      string `json:"name"`
	MTU       int    `json:"mtu"`
	RxPackets int    `json:"rx_pkt"`
	RxBytes   int    `json:"rx_bts"`
	TxPackets int    `json:"tx_pkt"`
	TxBytes   int    `json:"tx_bts"`
}
