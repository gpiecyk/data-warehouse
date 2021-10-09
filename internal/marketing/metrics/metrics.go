package metrics

import (
	"time"

	"github.com/gpiecyk/data-warehouse/internal/marketing/campaigns"
	"github.com/gpiecyk/data-warehouse/internal/marketing/datasources"
	"github.com/gpiecyk/data-warehouse/internal/marketing/files"
	"gorm.io/gorm"
)

type Metric struct {
	gorm.Model

	DatasourceIdFk uint                   `gorm:"not null;index:idx_datasource_date,priority:1"`
	Datasource     datasources.Datasource `gorm:"foreignKey:DatasourceIdFk"`

	CampaignIdFk uint               `gorm:"not null;index:idx_campaign"`
	Campaign     campaigns.Campaign `gorm:"foreignKey:CampaignIdFk"`

	FileIdFk uint       `gorm:"not null"`
	File     files.File `gorm:"foreignKey:FileIdFk"`

	Date        time.Time `json:"date,omitempty" gorm:"not null;index:idx_datasource_date,priority:2,sort:DESC;index:idx_date,sort:DESC"`
	Clicks      int       `json:"clicks,omitempty" gorm:"not null"`
	Impressions int       `json:"impressions,omitempty" gorm:"not null"`
}
