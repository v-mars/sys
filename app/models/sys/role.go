package sys

import (
	"github.com/v-mars/frame/db"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/sys/app/models"
	"github.com/v-mars/sys/app/utils"
	"gorm.io/gorm"
)

type Role struct {
	db.BaseID
	Title       string          `gorm:"type:varchar(128);unique_index:idx_title_code;not null;comment:'名称'" json:"title" form:"title"`
	Name        string          `gorm:"type:varchar(128);unique_index:idx_name_code;not null;comment:'名称'" json:"name" form:"name"`
	Description string          `gorm:"type:longtext;comment:'描述'" json:"description" form:"description"`
	ParentID    uint            `gorm:"index:parent_id;comment:'父节点'" json:"parent_id" form:"parent_id"`
	Path        models.IntArray `json:"path" gorm:"type:text;comment:'节点路径'"`
	Sort        int             `json:"sort" gorm:"default:1000;comment:'排序'"`
	Apis        []Api           `gorm:"many2many:role_api;association_autoupdate:false;association_autocreate:false;constraint:OnDelete:CASCADE" json:"apis"`
	//Children []Api `json:"children" gorm:"-"`
	db.BaseByUpdate
	db.BaseTime
}

func (Role) TableName() string {
	return "role"
}

func (d Role) GetMapTree(DB *gorm.DB) (map[uint]map[string]interface{},error){
	var data []Role
	if err:=DB.Order("parent_id asc,sort asc").Find(&data).Error;err!=nil{
		return nil, err
	}
	var MapData []map[string]interface{}
	var err error
	if MapData,err=convert.StructToMapSlice(data); err!=nil{
		return nil, err
	}
	var treeData  = map[uint]map[string]interface{}{}
	if len(MapData)>0{
		treeData= utils.ListToTree(MapData)
	}
	return treeData, nil
}

func (d Role) GetListTree(DB *gorm.DB) ([]map[string]interface{},error){
	treeMap, err := d.GetMapTree(DB)
	if err != nil {
		return nil, err
	}
	var dataList = treeMap[0]["children"].([]map[string]interface{})
	return dataList, nil
}