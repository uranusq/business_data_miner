package main

//Industries ... Industried by `Thomson Reuters Business Classification` and number of files belonging to them
type Industries struct {
	//gorm.Model
	ID       int    `gorm:"primary_key;AUTO_INCREMENT"`
	Industry string `gorm:"unique;not null"`
	NumURL   *uint  `gorm:"default:0"`
	NumHTML  *uint  `gorm:"default:0"`
	NumDocs  *uint  `gorm:"default:0"`
}

// Companies ... Companies with URL and other info that belong to some industry
type Companies struct {
	//gorm.Model
	ID              int    `gorm:"primary_key;AUTO_INCREMENT"`
	URL             string `gorm:"unique;not null"`
	Name            string
	IsCommonCrawled bool   `gorm:"default:0"`
	IsGoogleFilter  bool   `gorm:"default:0"`
	IsCollyCrawled  bool   `gorm:"default:0"`
	NumDocs         *uint  `gorm:"default:0"`
	NumHTML         *uint  `gorm:"default:0"`
	Industry        string `sql:"type:integer REFERENCES Industries(industry)"`
}
