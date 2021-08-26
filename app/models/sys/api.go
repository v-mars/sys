package sys

import (
	"github.com/v-mars/frame/db"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/sys/app/utils"
	"gorm.io/gorm"
)


var (
	CategoryApi = "api"
	Category = "category"
	CategoryMenu = "menu"
)

type Api struct {
	db.BaseID
	Name     string          `gorm:"type:varchar(128);not null;index:name;comment:名称" json:"name" form:"name"`
	Title    string          `json:"title" gorm:"size:128;index:title;comment:标题"`
	Group    string          `json:"group" gorm:"default:'default';size:128;comment:组"`
	Path     string          `gorm:"type:varchar(128);unique_index:idx_name_code;comment:路由路径" json:"path" form:"path"`
	Method   string          `gorm:"type:varchar(64);unique_index:idx_name_code;comment:请求方法" json:"method" form:"method"`
	Disabled bool            `json:"disabled" gorm:"default:0;comment:禁用"`
	//RoleId   uint   `gorm:"-"`
	//ParentID uint            `gorm:"index:parent_id;comment:'父节点'" json:"parent_id" form:"parent_id"`
	//Paths    models.IntArray `json:"paths" gorm:"type:text;comment:'节点路径'"`
	//Sort     int             `json:"sort" gorm:"default:1000;comment:'排序'"`
	//Category string          `json:"category" gorm:"default:'api';not null;index:category;size:16;comment:'分类：api,category,menu';"`
	//Children []Api `json:"children" gorm:"-"`
	db.BaseTime
}

func (Api) TableName() string {
	return "api"
}


func (d Api) GetMapTree(DB *gorm.DB) (map[uint]map[string]interface{},error){
	var data []Api
	if err:=DB.Order("parent_id asc,path asc").Find(&data).Error;err!=nil{
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

func (d Api) GetListTree(DB *gorm.DB) ([]map[string]interface{},error){
	treeMap, err := d.GetMapTree(DB)
	if err != nil {
		return nil, err
	}
	var dataList = treeMap[0]["children"].([]map[string]interface{})
	return dataList, nil
}