package property

import (
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/frame/pkg/utils"
	"github.com/v-mars/frame/response"
)

type Gitlab struct {
	Token  string `json:"token"`
	Host   string `json:"host"`
}
func (r Property) QueryGitlab(c *gin.Context) {
	data,err := GetProperty(r.DB,"gitlab")
	if err != nil{
		response.Error(c, err)
		return
	}
	if data["token"] != "" {
		data["token"] = utils.Base64Dec(convert.ToString(data["token"]))
	}
	response.Success(c, map[string]interface{}{"result": data})
}

func (r Property) CreateOrUpdateGitlab(c *gin.Context) {
	type Param struct {
		Rows        map[string]interface{} `json:"row" binding:"required"`
	}
	var obj Param
	if err := c.ShouldBindJSON(&obj); err!=nil{
		response.ParamFailed(c, err)
		return
	}
	obj.Rows["token"] = utils.Base64Enc([]byte(convert.ToString(obj.Rows["token"])))
	//fmt.Println("obj.Rows:", obj.Rows)
	if err:= CreateOrUpdate(r.DB,"gitlab", obj.Rows); err!=nil{
		response.Error(c, err)
		return
	}
	response.SuccessMsg(c, "保存成功！", "")
}
