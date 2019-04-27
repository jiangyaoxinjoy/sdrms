package controllers

import (
	"encoding/json"
	"fmt"
	"sdrms/enums"
	"sdrms/models"
	"sdrms/utils"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type ApiController struct {
	BaseController
}

//获取所有文章 limit和page
func (c *ApiController) CourseGrid() {
	params := models.ApiCourseQueryParam{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &params); err != nil {
		c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("%s", err), "")
	}
	if params.Limit == 0 {
		params.Limit = 5
	}
	offsetint := (params.Page - 1) * params.Limit
	fmt.Println(offsetint)
	offsetstr := strconv.Itoa(offsetint)
	offset64, _ := strconv.ParseInt(offsetstr, 10, 64)
	params.Offset = offset64

	data, err := models.ApiCoursePageList(&params)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "", data)
	}

	postData := make([]models.ApiCourse, len(data))
	for i := 0; i < len(data); i++ {
		fmt.Println("*****************")
		fmt.Println(data[i].Comments)
		postData[i].Title = data[i].Title
		postData[i].CreatedTime = data[i].CreatedTime
		postData[i].Identify = data[i].Identify
		postData[i].Id = data[i].Id
		postData[i].Identify = data[i].Identify
		postData[i].Status = data[i].Status
		postData[i].Content = data[i].CourseContent.Content
		postData[i].CommentNum = len(data[i].Comments)
		for _, val := range data[i].Tags {
			postData[i].Tags = append(postData[i].Tags, val.Name)
		}
	}
	c.jsonResult(enums.JRCodeSucc, "ok", postData)
}

func (c *ApiController) OneCourse() {
	id, _ := c.GetInt(":aid")
	if id < 0 {
		c.jsonResult(enums.JRCodeFailed, "数据不存在", "")
	}
	var data models.ApiCourse
	course, err := models.ApiGetCourseById(id)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("%s", err), "")
	}
	data.Title = course.Title
	data.CreatedTime = course.CreatedTime
	data.Identify = course.Identify
	data.Id = course.Id
	data.Identify = course.Identify
	data.Status = course.Status
	data.Content = course.CourseContent.Content
	data.CommentNum = len(course.Comments)
	for _, val := range course.Tags {
		data.Tags = append(data.Tags, val.Name)
	}
	c.jsonResult(enums.JRCodeSucc, "ok", data)
}

func (c *ApiController) Categorys() {
	if c.Ctx.Request.Method == "POST" {
		c.getCategoryById()
	}
	data, err := models.ApiGetCategorys()
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据不存在", "")
	}
	c.jsonResult(enums.JRCodeSucc, "ok", data)
}

func (c *ApiController) getCategoryById() {
	if id, err := c.GetInt(":id"); err == nil {
		if id > 0 {
			if data, err := models.ApiGetOneCategoryById(id); err == nil {
				c.jsonResult(enums.JRCodeSucc, "ok", data)
			} else {
				c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("%s", err), "")
			}
		}
	} else {
		c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("%s", err), "")
	}
}

func (c *ApiController) Login() {
	user := new(models.User)
	json.Unmarshal(c.Ctx.Input.RequestBody, user)
	if len(user.Username) == 0 || len(user.Userpwd) == 0 {
		c.jsonResult(enums.JRCodeFailed, "用户名和密码不正确", "")
	}
	user.Userpwd = utils.String2md5(user.Userpwd)
	if userMsg, err := models.BackendUserOneByUserName(user.Username, user.Userpwd); userMsg != nil && err == nil {
		c.setBackendUser2Session(userMsg.Id)
		apiUser := new(models.ApiUser)
		apiUser.Name = userMsg.UserName
		token := CreateToken(apiUser)
		apiUser.Token = token
		c.jsonResult(enums.JRCodeSucc, "登录成功", apiUser)
	} else {
		c.jsonResult(enums.JRCodeFailed, "用户名或者密码错误", "")
	}
}

func CreateToken(user *models.ApiUser) string {
	type UserInfo map[string]interface{}
	key := fmt.Sprintf("welcome to %s's code world", user.Name)
	userInfo := make(UserInfo)
	userInfo["exp"] = "1515482650719371100" //  strconv.FormatInt(t.UTC().UnixNano(), 10)
	userInfo["iat"] = "0"

	tokenString := models.ApiCreateToken(key, userInfo)
	return tokenString
}

func (c *ApiController) ParseToken() {
	t := time.Now()
	// tokenString := c.Ctx.GetCookie("token")
	tokenString := c.Ctx.Input.Header("Authorization")
	key := fmt.Sprintf("welcome to %s's code world", c.curUser.UserName)
	claims, ok := models.ApiParseToken(tokenString, key)
	var tokenState string
	var expTime int64 = 10
	if ok {
		oldT, _ := strconv.ParseInt(claims.(jwt.MapClaims)["exp"].(string), 10, 64)
		ct := t.UTC().UnixNano()
		cur := ct - oldT
		fmt.Println("************5555555555555")
		fmt.Println(cur)
		if cur > expTime {
			ok = false
			tokenState = "Token 已过期"
			c.jsonResult(enums.JRCodeFailed, tokenState, "")
		} else {
			tokenState = "Token 正常"
			c.jsonResult(enums.JRCodeSucc, tokenState, "")
		}
	} else {
		tokenState = "Token 无效"
		c.jsonResult(enums.JRCodeFailed, tokenState, "")
	}
}

//获取所有标签
func (c *ApiController) GetTags() {
	tags := models.AllCategorys()
	data := make([]models.ApiTags, len(tags))
	for i, val := range tags {
		data[i].Name = val.Name
		data[i].Id = val.Id
	}
	c.jsonResult(enums.JRCodeSucc, "ok", data)
}

func (c *ApiController) GetComments() {
	if c.Ctx.Request.Method == "POST" {
		c.saveComment()
	}
	identify := c.GetString("id")
	fmt.Println(identify)
	comments, err := models.ApiGetCommentByIdentify(identify)
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "数据无效", "")
	}
	data := make([]models.ApiComment, len(comments))
	for i, val := range comments {
		data[i].NickName = val.NickName
		data[i].CreatedTime = val.CreatedTime
		data[i].ImgName = val.ImgName
		data[i].Content = val.Content
		data[i].Like = val.Like
	}
	c.jsonResult(enums.JRCodeSucc, "ok", data)
}

func (c *ApiController) saveComment() {
	data := new(models.ApiComment)
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, data); err != nil {
		c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("%s", err), "")
	}
	ctime := utils.GetCurrTs()
	data.CreatedTime = ctime
	if err := models.ApiAddComment(data); err != nil {
		c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("%s", err), "")
	}
	c.jsonResult(enums.JRCodeSucc, "ok", "")
}
