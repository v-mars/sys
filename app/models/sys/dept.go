package sys

import (
	"github.com/v-mars/frame/db"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/sys/app/models"
	"github.com/v-mars/sys/app/utils"
	"gorm.io/gorm"
)

type Dept struct {
	db.BaseID
	Name     string `gorm:"type:varchar(156);not null;unique:idx_name;comment:'名称'" json:"name" form:"name"`
	Code     string `gorm:"type:varchar(128);not null;unique:idx_code;comment:'代号'" json:"code" form:"code"`
	ParentID uint   `gorm:"index:parent_id;comment:'父节点'" json:"parent_id" form:"parent_id"`
	Path     models.IntArray  `json:"path" gorm:"type:text;comment:'节点路径'"`
	Sort     int    `json:"sort" gorm:"comment:'排序'"`
	Leader   string `json:"leader" gorm:"size:128;comment:'负责人'"`                 //负责人
	Status   bool   `json:"status" gorm:"default:1;comment:'状态'"`                   //状态
	Children []Dept `json:"children" gorm:"-"`
	db.BaseTime
}

func (Dept) TableName() string {
	return "dept"
}

func (d Dept) GetMapTree(DB *gorm.DB) (map[uint]map[string]interface{},error){
	var data []Dept
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

func (d Dept) GetListTree(DB *gorm.DB) ([]map[string]interface{},error){
	treeMap, err := d.GetMapTree(DB)
	if err != nil {
		return nil, err
	}
	var dataList = treeMap[0]["children"].([]map[string]interface{})
	return dataList, nil
}


