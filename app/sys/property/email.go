package property

import (
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/frame/pkg/utils"
	"github.com/v-mars/frame/response"
)

type Email struct {
	Token  string `json:"token"`
	Host   string `json:"host"`
}
func (r Property) QueryEmail(c *gin.Context) {
	data,err := GetProperty(r.DB,"email")
	if err != nil{
		response.Error(c, err)
		return
	}
	if data["password"] != "" {
		data["password"] = utils.Base64Dec(convert.ToString(data["password"]))
	}
	response.Success(c, map[string]interface{}{"result": data})
}

func (r Property) CreateOrUpdateEmail(c *gin.Context) {
	type Param struct {
		//CheckEmail  string                 `json:"check_email"`
		Rows        map[string]interface{} `json:"row" binding:"required"`
	}
	var obj Param
	if err := c.ShouldBindJSON(&obj); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	obj.Rows["password"] = utils.Base64Enc([]byte(convert.ToString(obj.Rows["password"])))
	//fmt.Println("obj.Rows:", obj.Rows)
	if err:= CreateOrUpdate(r.DB,"email", obj.Rows); err!=nil{
		response.Error(c, err)
		return
	}
	response.SuccessMsg(c, "保存成功！", "")
}