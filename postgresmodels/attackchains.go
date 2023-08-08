package postgresmodels

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

/*
Related resources, controls info and posture data will be enriched in backend business logic level until ready in postgres.
*/

type AttackChainState struct {
	// BaseModel.CreatedAt is the former FirstSeen and CreationTime which are the same
	BaseModel

	// primary keys
	AttackChainID   string `gorm:"primaryKey;not null"` // name/cluster/resourceID
	CustomerGUID    string `gorm:"primaryKey;not null"`
	AttackTrackName string `gorm:"primaryKey;not null"`

	AttackTrackDescription string

	// attributes["cluster"], attributes["namespace"], attributes["kind"], attributes["name"]
	Resource           datatypes.JSON
	ResourceName       string `gorm:"not null"`
	ResourceNamespace  string `gorm:"not null"`
	ResourceKind       string `gorm:"not null"`
	ResourceAPIVersion string `gorm:"not null"`
	ClusterName        string `gorm:"not null"`

	LatestReportGUID string `gorm:"not null"` // latest reportGUID in which this attack chain was identified

	Status string // "active"/ "fixed"

	// processing status is updated by the UI once a scan is initiated for all relevant clusters (connected) of the customerGUID.
	// "done" is updated by the attack chain engine once finished processing.
	ProcessingStatus string    `gorm:"not null"` // "processing"/ "done"
	ViewedMainScreen time.Time // updated by UI - if the attack chain was viewed by the user// New badge

	RootNode   AttackChainNode `gorm:"foreignKey:RootNodeID"`
	RootNodeID uint            `gorm:"not null"`
}

type AttackChainNode struct {
	gorm.Model           // ID, CreatedAt, UpdatedAt, DeletedAt - ID is required for linking nodes
	Name          string `gorm:"not null"`
	AttackChainID string `gorm:"not null"` // hash of cluster/resourceID
	CustomerGUID  string `gorm:"not null"`
	IsRoot        bool   `gorm:"not null"`
}

type AttackChainNodeRelation struct {
	BaseModel
	ParentNode   AttackChainNode `gorm:"foreignKey:ParentNodeID"`
	ParentNodeID uint            `gorm:"primaryKey; not null"`
	ChildNode    AttackChainNode `gorm:"foreignKey:ChildNodeID"`
	ChildNodeID  uint            `gorm:"primaryKey; not null"`
}

type AttackChainNodeImageScanRelation struct {
	BaseModel
	NodeID uint            `gorm:"primaryKey; not null"`
	Node   AttackChainNode `gorm:"foreignKey:NodeID"`

	// ImageScanId = hash of customerGUID, cluster, containerSpecID
	// Should be used instead of ContainersScanID
	ImageScanId string `gorm:"primaryKey; not null"`

	// TODO: define ImageScanSummary with foreign key - need to fix TestVulScan dumb data tests in postgres connector to be aligned with key constaints
	// ImageScanSummary VulnerabilityScanSummary `gorm:"foreignKey:ImageScanId"`
}

type AttackChainNodeRelatedResourcesRelation struct {
	BaseModel
	NodeID     uint            `gorm:"primaryKey; not null"`
	Node       AttackChainNode `gorm:"foreignKey:NodeID"`
	ResourceID string          `gorm:"primaryKey; not null"`
}

type AttackChainNodeControlsRelation struct {
	BaseModel
	NodeID uint            `gorm:"primaryKey; not null"`
	Node   AttackChainNode `gorm:"foreignKey:NodeID"`

	// ControlID = failed or ignored control ID that is associated with the node.
	ControlID string `gorm:"primaryKey; type:varchar(255);not null"`
}

func (AttackChainState) TableName() string {
	return "attack_chains_states"
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
