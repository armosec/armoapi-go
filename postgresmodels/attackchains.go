package postgresmodels

import (
	"gorm.io/gorm"
)

type AttackChainNode struct {
	gorm.Model            // ID, CreatedAt, UpdatedAt, DeletedAt - ID is required for linking nodes
	Name          string  `gorm:"type:varchar(255);not null"`
	Description   *string `gorm:"type:varchar(255)" `
	AttackChainID string  `gorm:"primaryKey;type:varchar(255);not null"`
	CustomerGUID  string  `gorm:"primaryKey;type:varchar(255);not null"`
}

type AttackChainNodeRelation struct {
	BaseModel
	ParentNodeID uint `gorm:"primaryKey; not null"`
	ChildNodeID  uint `gorm:"primaryKey; not null"`
}

type AttackChainNodeImageScanRelation struct {
	BaseModel
	NodeID uint `gorm:"primaryKey; not null"`

	// ImageScanId = ContainerScanId (required for attack chain.)
	ImageScanId string `gorm:"primaryKey; foreignKey:ImageScanId; type:varchar(255);not null"`
}

type AttackChainNodeRelatedResourcesRelation struct {
	BaseModel
	NodeID     uint   `gorm:"primaryKey; not null"`
	ResourceID string `gorm:"primaryKey; type:varchar(255);not null"`
}

type AttackChainNodeControlsRelation struct {
	BaseModel
	NodeID uint `gorm:"primaryKey; not null"`

	// ControlID = failed or ignored control ID that is associated with the node.
	ControlID string `gorm:"primaryKey; type:varchar(255);not null"`
}

func (AttackChainNode) TableName() string {
	return "attack_chain_nodes"
}

func (AttackChainNodeRelation) TableName() string {
	return "attack_chain_nodes_relations"
}

func (AttackChainNodeImageScanRelation) TableName() string {
	return "attack_chain_node_image_scan_relations"
}

func (AttackChainNodeRelatedResourcesRelation) TableName() string {
	return "attack_chain_node_related_resources_relations"
}

func (AttackChainNodeControlsRelation) TableName() string {
	return "attack_chain_node_controls_relations"
}

// foreign keys and indexes should be manually executed:
/*
db.Exec("ALTER TABLE attack_chain_node_relations ADD CONSTRAINT fk_parent_node FOREIGN KEY (parent_node_id) REFERENCES attack_chain_nodes(id)")
db.Exec("ALTER TABLE attack_chain_node_relations ADD CONSTRAINT fk_child_node FOREIGN KEY (child_node_id) REFERENCES attack_chain_nodes(id)")
db.Exec("CREATE INDEX idx_attack_chain_node_relations_parent_child ON attack_chain_node_relations (parent_node_id, child_node_id);")
*/
