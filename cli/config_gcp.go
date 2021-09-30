package cli

type GCPConfig struct {
	GCPOrgID         string `json:",omitempty"`
	GCPBillingID     string `json:",omitempty"`
	GCPParentProject string `json:",omitempty"`
	GCPUserEmail     string `json:",omitempty"`
	GCPDNSZone       string `json:",omitempty"`
}

func (g *GCPConfig) SetGCPOrgID(id string) {
	g.GCPOrgID = "organizations/" + id
}

func (g *GCPConfig) SetGCPBillingID(id string) {
	g.GCPBillingID = id
}
