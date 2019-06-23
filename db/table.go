package db

//Industries ... Industried by `Thomson Reuters Business Classification` and number of files belonging to them
type Industries struct {
	//gorm.Model
	ID       int    `gorm:"primary_key;AUTO_INCREMENT"`
	Industry string `gorm:"unique;not null"`
	NumURL   *uint  `gorm:"default:0"`
	NumHTML  *uint  `gorm:"default:0"`
	NumDocs  *uint  `gorm:"default:0"`
}

//Industries ... Industried by `Thomson Reuters Business Classification` and number of files belonging to them
type IndustryGroups struct {
	//gorm.Model
	ID             int    `gorm:"primary_key;AUTO_INCREMENT"`
	IndustryGroups string `gorm:"unique;not null"`
	NumURL         *uint  `gorm:"default:0"`
	NumHTML        *uint  `gorm:"default:0"`
	NumDocs        *uint  `gorm:"default:0"`
}

//Industries ... Industried by `Thomson Reuters Business Classification` and number of files belonging to them
type Business struct {
	//gorm.Model
	ID       int    `gorm:"primary_key;AUTO_INCREMENT"`
	Business string `gorm:"unique;not null"`
	NumURL   *uint  `gorm:"default:0"`
	NumHTML  *uint  `gorm:"default:0"`
	NumDocs  *uint  `gorm:"default:0"`
}

//Industries ... Industried by `Thomson Reuters Business Classification` and number of files belonging to them
type Economics struct {
	//gorm.Model
	ID        int    `gorm:"primary_key;AUTO_INCREMENT"`
	Economics string `gorm:"unique;not null"`
	NumURL    *uint  `gorm:"default:0"`
	NumHTML   *uint  `gorm:"default:0"`
	NumDocs   *uint  `gorm:"default:0"`
}

// Companies ... Companies with URL and other info that belong to some industry
type Companies struct {
	//gorm.Model
	ID              int    `gorm:"primary_key;AUTO_INCREMENT"`
	URL             string `gorm:"unique;not null"`
	Name            string
	IsCommonCrawled bool   `gorm:"default:0"`
	IsGoogleCrawled bool   `gorm:"default:0"`
	IsCollyCrawled  bool   `gorm:"default:0"`
	NumDocs         *uint  `gorm:"default:0"`
	NumHTML         *uint  `gorm:"default:0"`
	Industry        string `sql:"type:integer REFERENCES Industries(industry)"`
	IndustryGroup   string `sql:"type:integer REFERENCES IndustryGroups(industry_group)"`
	Business        string `sql:"type:integer REFERENCES Business(business)"`
	Economics       string `sql:"type:integer REFERENCES Economics(economics)"`
}
