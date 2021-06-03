package property

import (
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/db"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/frame/response"
	"github.com/v-mars/sys/app/models"
	"github.com/v-mars/sys/app/models/sys"
	"gorm.io/gorm"
)

type IProperty interface {
	Query(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	QueryLDAP(c *gin.Context)
	CreateOrUpdateLDAP(c *gin.Context)
	CheckLDAPConnect(c *gin.Context)
	QueryEmail(c *gin.Context)
	CreateOrUpdateEmail(c *gin.Context)
	QueryAliYun(c *gin.Context)
	CreateOrUpdateAliYun(c *gin.Context)
	QueryGitlab(c *gin.Context)
	CreateOrUpdateGitlab(c *gin.Context)
	QueryJenkins(c *gin.Context)
	CreateOrUpdateJenkins(c *gin.Context)
}

type Property struct {
	DB *gorm.DB
}

func NewService(DB *gorm.DB) IProperty {
	if DB==nil{DB=models.DB}
	return Property{DB: DB}
}

func (r Property) Query(c *gin.Context) {
	var obj []sys.Property
	var pageData = response.InitPageData(c, &obj, false)
	type Param struct {
		Name string `form:"name"`
	}
	var param Param
	if err := c.ShouldBindQuery(&param);err!=nil{
		response.Error(c, err)
		return
	}
	var o db.Option
	o.Where = "name like ?"
	o.Value = append(o.Value, "%"+param.Name+"%")
	o.Select = "distinct id, name, k, v, ctime, mtime"
	o.Order = "ID DESC,name DESC"
	o.Scan = true

	err := db.Query(r.DB,&sys.Property{}, o, &pageData)
	if err !=nil {
		response.Error(c, err)
		return
	}else {
		response.Success(c, pageData)
		return
	}
}

func (r Property) Create(c *gin.Context) {

}

func (r Property) Update(c *gin.Context) {

}

func (r Property) Delete(c *gin.Context) {
	var obj sys.Property
	var rows map[string][]int
	if err := c.ShouldBindJSON(&rows); err!=nil{
		response.Error(c, err)
		return
	}
	if err:= r.DB.Where("id in (?)", rows["rows"]).Delete(&obj).Error;err!=nil{
		response.Error(c, err)
		return
	}
	response.DeleteSuccess(c)
}


func GetProperty(DB *gorm.DB,name string) (map[string]interface{}, error) {
	var o db.Option
	o.Where = "name = ?"
	o.Value = append(o.Value, name)
	o.Select = "distinct id, name, k, v, ctime, mtime"
	o.Order = "ID DESC"
	o.Scan = true
	var obj []sys.Property
	var pageData = response.PageDataList{PageNumber: 1,PageSize:0,List:&obj}
	err := db.Query(DB,&sys.Property{}, o, &pageData)
	if err!=nil{
		return nil, err
	}
	var data = map[string]interface{}{}
	for _, v := range obj {
		data[v.K] = v.V
		//fmt.Println("kv:",  v, v.K)
	}
	return data, nil
}

func CreateOrUpdate(DB *gorm.DB,name string,m map[string]interface{}) error {
	for k, v := range m {
		if err:= DB.Where(sys.Property{Name: name, K: k}).Assign(
			sys.Property{Name: name, K: k, V: convert.ToString(v)},
			).FirstOrCreate(&sys.Property{}).Error; err!=nil{
				return err
		}
	}
	return nil
}