package postgresmodels

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AttackChainNode struct {
	gorm.Model                      // ID, CreatedAt, UpdatedAt, DeletedAt - ID is required for linking nodes
	Name             string         `gorm:"type:varchar(255);not null" json:"name"`
	Description      *string        `gorm:"type:varchar(255)" json:"description"`
	ControlIDs       datatypes.JSON `gorm:"type:jsonb" json:"controlIDs,omitempty"`
	Vulnerabilities  datatypes.JSON `gorm:"type:jsonb" json:"vulnerabilitiesNames,omitempty"`
	RelatedResources datatypes.JSON `gorm:"type:jsonb" json:"relatedResources"`
	AttackChainID    string         `gorm:"type:varchar(255);not null" json:"attackChainID"`
	CustomerGUID     string         `gorm:"type:varchar(255);not null" json:"customerGUID"`
}

type AttackChainNodeLink struct {
	BaseModel
	ParentNodeID uint `gorm:"not null"`
	ChildNodeID  uint `gorm:"not null"`
}

func (AttackChainNode) TableName() string {
	return "attack_chain_nodes"
}

func (AttackChainNodeLink) TableName() string {
	return "attack_chain_nodes_links"
}

// foreign keys and indexes should be manually executed:
/*
db.Exec("ALTER TABLE attack_chain_node_links ADD CONSTRAINT fk_parent_node FOREIGN KEY (parent_node_id) REFERENCES attack_chain_nodes(id)")
db.Exec("ALTER TABLE attack_chain_node_links ADD CONSTRAINT fk_child_node FOREIGN KEY (child_node_id) REFERENCES attack_chain_nodes(id)")
db.Exec("CREATE INDEX idx_attack_chain_node_links_parent_child ON attack_chain_node_links (parent_node_id, child_node_id);")
*/
