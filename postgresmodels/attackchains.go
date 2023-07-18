package postgresmodels

import (
	"gorm.io/gorm"
)

type AttackChainNode struct {
	gorm.Model            // ID, CreatedAt, UpdatedAt, DeletedAt - ID is required for linking nodes
	Name          string  `gorm:"type:varchar(255);not null" json:"name"`
	Description   *string `gorm:"type:varchar(255)" json:"description"`
	AttackChainID string  `gorm:"type:varchar(255);not null" json:"attackChainID"`
	CustomerGUID  string  `gorm:"type:varchar(255);not null" json:"customerGUID"`
}

type AttackChainNodeRelation struct {
	BaseModel
	ParentNodeID uint `gorm:"not null"`
	ChildNodeID  uint `gorm:"not null"`
}

type AttackChainNodeImageScanRelations struct {
	BaseModel
	NodeID uint `gorm:"not null"`

	// via VulnerabilityScanSummary.ImageScanId, can get the ContainerSpecId (ContainerScanId) and list of image names required for attack chain.
	ImageScanId string `gorm:"type:varchar(255);not null" json:"imageScanId"`
}

type AttackChainNodeRelatedResourcesRelations struct {
	BaseModel
	NodeID     uint   `gorm:"not null"`
	ResourceID string `gorm:"type:varchar(255);not null" json:"resourceID"`
}

type AttackChainNodeControlsRelations struct {
	BaseModel
	NodeID    uint   `gorm:"not null"`
	ControlID string `gorm:"type:varchar(255);not null" json:"controlID"`
}

func (AttackChainNode) TableName() string {
	return "attack_chain_nodes"
}

func (AttackChainNodeRelation) TableName() string {
	return "attack_chain_nodes_relations"
}

func (AttackChainNodeImageScanRelations) TableName() string {
	return "attack_chain_node_image_scan_relations"
}

func (AttackChainNodeRelatedResourcesRelations) TableName() string {
	return "attack_chain_node_related_resources_relations"
}

func (AttackChainNodeControlsRelations) TableName() string {
	return "attack_chain_node_controls_relations"
}

// foreign keys and indexes should be manually executed:
/*
db.Exec("ALTER TABLE attack_chain_node_relations ADD CONSTRAINT fk_parent_node FOREIGN KEY (parent_node_id) REFERENCES attack_chain_nodes(id)")
db.Exec("ALTER TABLE attack_chain_node_relations ADD CONSTRAINT fk_child_node FOREIGN KEY (child_node_id) REFERENCES attack_chain_nodes(id)")
db.Exec("CREATE INDEX idx_attack_chain_node_relations_parent_child ON attack_chain_node_relations (parent_node_id, child_node_id);")
*/
