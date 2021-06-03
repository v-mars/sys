package property

import (
	"github.com/gin-gonic/gin"
	"github.com/goinggo/mapstructure"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/frame/pkg/utils"
	"github.com/v-mars/frame/response"
	"gorm.io/gorm"
)

type Jenkins struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}
func (r Property) QueryJenkins(c *gin.Context) {
	data,err := GetProperty(r.DB,"jenkins")
	if err != nil{
		response.Error(c, err)
		return
	}
	if data["password"] != "" {
		data["password"] = utils.Base64Dec(convert.ToString(data["password"]))
	}
	response.Success(c, map[string]interface{}{"result": data})
}

func (r Property) CreateOrUpdateJenkins(c *gin.Context) {
	type Param struct {
		Rows        map[string]interface{} `json:"row" binding:"required"`
	}
	var obj Param
	if err := c.ShouldBindJSON(&obj); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	obj.Rows["password"] = utils.Base64Enc([]byte(convert.ToString(obj.Rows["password"])))
	//fmt.Println("obj.Rows:", obj.Rows)
	if err:= CreateOrUpdate(r.DB,"jenkins", obj.Rows); err!=nil{
		response.Error(c, err)
		return
	}
	response.SuccessMsg(c, "保存成功！", "")
}


func GetJenkins(DB *gorm.DB) (Jenkins, error) {
	data,err := GetProperty(DB,"jenkins")
	var j Jenkins
	if err != nil{
		return j, err
	}
	if data["password"] != "" {
		data["password"] = utils.Base64Dec(convert.ToString(data["password"]))
	}
	err = mapstructure.Decode(data, &j)
	if err != nil {
		return j, err
	}
	return j, err
}
