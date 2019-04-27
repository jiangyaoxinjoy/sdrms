package controllers

import (
	"encoding/json"
	"fmt"
	"sdrms/enums"
	"sdrms/models"
	"sdrms/utils"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type CourseController struct {
	BaseController
}

func (c *CourseController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的少数Action需要权限控制，则将验证放到需要控制的Action里
	//c.checkAuthor("TreeGrid", "UserMenuTree", "ParentTreeGrid", "Select")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//这里注释了权限控制，因此这里需要登录验证
	c.checkLogin()
}

func (c *CourseController) Index() {
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "course/index_footerjs.html"
	c.LayoutSections["headcssjs"] = "course/index_headcssjs.html"
	c.LayoutSections["pagejs"] = "shared/paging.html"
	c.Data["canEdit"] = c.checkActionAuthor("CourseController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("CourseController", "Delete")
	c.Data["showMoreQuery"] = true
	category := models.AllCategorys()
	c.Data["category"] = category
}

func (c *CourseController) DataGrid() {
	var params models.CourseQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//获取数据列表和总数
	data, total := models.CoursePageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *CourseController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.save()

	}
	Id, _ := c.GetInt(":id", 0)
	if Id == 0 {

	} else {
		course, err := models.CourseOne(Id)
		if err != nil {
			c.pageError(fmt.Sprintf("%s", err))
		}
		c.Data["course"] = course
	}
	category := models.AllCategorys()
	allTags := models.AllTags()
	c.Data["category"] = category
	c.Data["allTags"] = allTags
	c.setTpl("course/edit.html", "shared/layout_base.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "course/edit_footerjs.html"
	c.LayoutSections["headcssjs"] = "course/edit_headcssjs.html"
}

func (c *CourseController) save() {
	course := &models.Course{}

	o := orm.NewOrm()
	o.Begin()

	id, _ := c.GetInt("Id", 0)
	tags := c.GetString("Tags")
	tagsarrs := strings.Split(tags, ",")
	tagsarr := make([]int, len(tagsarrs))
	for _, val := range tagsarrs {
		if s, errs := strconv.Atoi(val); errs == nil {
			tagsarr = append(tagsarr, s)
		}
	}
	courseCatetoryId, _ := c.GetInt("CourseCatetory")
	content := c.GetString("content-markdown-doc")
	title := c.GetString("Title")
	status, _ := c.GetInt("Status", 0)

	valid := validation.Validation{}
	valid.Required(title, "标题")
	valid.Required(tags, "关键词")
	valid.Required(content, "内容")
	valid.Required(courseCatetoryId, "分类")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.jsonResult(enums.JRCodeFailed, "错误提示", fmt.Sprintf("%s %s", err.Name, err.Message))
		}
	}
	ctime := utils.GetCurrTs()
	if id > 0 {
		course.Id = id
		if err := o.Read(course); err != nil {
			c.jsonResult(enums.JRCodeFailed, "数据获取失败", fmt.Sprintf("%s", err))
		}
		m2m := o.QueryM2M(course, "Tags")
		m2m.Clear()
		tags := make([]*models.Tag, len(tags))
		models.GetTagIds(tagsarr, &tags)
		m2m.Add(tags)

		category := &models.CourseCategory{Id: courseCatetoryId}
		course.CourseCategory = category

		course.Title = title
		course.Status = status

		if _, err := o.QueryTable(models.CourseContentTBName()).Filter("Course__Id", id).Update(orm.Params{"content": content}); err != nil {
			o.Rollback()
			c.jsonResult(enums.JRCodeFailed, "修改失败", fmt.Sprintf("%s", err))
		}
		course.ModifyTime = ctime
		if _, err := o.Update(course, "Title", "Status", "CourseCategory", "Tags", "ModifyTime"); err != nil {
			o.Rollback()
			c.jsonResult(enums.JRCodeFailed, "修改失败", fmt.Sprintf("%s", err))
		}
		o.Commit()
		c.jsonResult(enums.JRCodeSucc, "ok", "")
	} else {
		course.CreatedTime = ctime
		course.ModifyTime = ctime
		num := utils.Generate_Randnum(20, 10)
		identify := utils.GetRandomString(num)
		course.Identify = identify
		course.Title = title
		course.Status = status
		backendUser := &models.BackendUser{Id: c.curUser.Id}
		course.BackendUser = backendUser

		coursecontent := &models.CourseContent{Content: content}
		if cid, err := o.Insert(coursecontent); err == nil {
			strInt64 := strconv.FormatInt(cid, 10)
			id16, _ := strconv.Atoi(strInt64)
			coursecontent.Id = id16
			course.CourseContent = coursecontent
		}

		category := &models.CourseCategory{Id: courseCatetoryId}
		course.CourseCategory = category

		courseid, _ := o.Insert(course)
		cids := strconv.FormatInt(courseid, 10)
		cidint, _ := strconv.Atoi(cids)
		course.Id = cidint
		m2m := o.QueryM2M(course, "Tags")
		tags := make([]*models.Tag, len(tags))
		if err := models.GetTagIds(tagsarr, &tags); err != nil {
			o.Rollback()
			c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("%s", err), "")
		}
		fmt.Println(tags)
		m2m.Add(tags)
		o.Commit()
		c.jsonResult(enums.JRCodeSucc, "ok", "")
	}
}

func (c *CourseController) EditormdPic() {
	f, h, err := c.GetFile("editormd-image-file")
	if err != nil {
		utils.LogError(fmt.Sprintf("getfile err: %d ", err))
	}

	defer f.Close()
	urlname := fmt.Sprintf("/static/upload/editormdImg/%s", h.Filename)
	if err := c.SaveToFile("editormd-image-file", "."+urlname); err != nil {
		utils.LogError(fmt.Sprintf("savefile err: %d ", err))
		c.Data["json"] = map[string]interface{}{"success": 0, "message": "图片上传失败", "url": ""}
		c.ServeJSON()
		c.StopRun()
	}
	c.Data["json"] = map[string]interface{}{"success": 1, "message": h.Filename, "url": urlname}
	c.ServeJSON()
	c.StopRun()

}

//批量删除课程
func (c *CourseController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	if num, err := models.CourseBatchDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}

func (c *CourseController) Tags() {
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl("tags/index.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "tags/index_footerjs.html"
	c.LayoutSections["headcssjs"] = "tags/index_headcssjs.html"
	c.Data["canEdit"] = c.checkActionAuthor("CourseController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("CourseController", "Delete")
}

func (c *CourseController) TagsDataGrid() {
	params := new(models.BaseQueryParam)
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	pagelist, total := models.TagsPageList(params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = pagelist
	c.Data["json"] = result
	c.ServeJSON()
}
func (c *CourseController) EidtTag() {
	if c.Ctx.Request.Method == "POST" {
		c.saveTag()
	}
	c.setTpl("tags/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "tags/edit_footerjs.html"
	c.LayoutSections["headcssjs"] = "tags/edit_headcssjs.html"
}

func (c *CourseController) saveTag() {
	tag := models.Tag{}
	if err := c.ParseForm(&tag); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", tag.Id)
	}
	if tag.Id > 0 {
		o := orm.NewOrm()
		if _, err := o.Update(&tag); err != nil {
			c.jsonResult(enums.JRCodeFailed, "数据插入失败", tag.Id)
		}
		c.jsonResult(enums.JRCodeSucc, "数据修改成功", tag.Id)
	} else {
		o := orm.NewOrm()
		if _, err := o.Insert(&tag); err != nil {
			c.jsonResult(enums.JRCodeFailed, fmt.Sprintf("%s", err), tag.Id)
		}
		c.jsonResult(enums.JRCodeSucc, "添加标签成功", tag.Id)
	}

}

func (c *CourseController) DeleteTag() {
	strs := c.GetString("ids")
	ids := make([]int, len(strs))
	for _, val := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(val); err == nil {
			ids = append(ids, id)
		}
	}
	if num, err := models.TagBatchDelete(ids); err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	} else {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	}
}

func (c *CourseController) Category() {
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl("coursecategory/index.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "coursecategory/index_footerjs.html"
	c.LayoutSections["headcssjs"] = "coursecategory/index_headcssjs.html"
	c.Data["canEdit"] = c.checkActionAuthor("CourseController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("CourseController", "Delete")
}

func (c *CourseController) CategoryDataGrid() {
	var params models.BaseQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//获取数据列表和总数
	data, total := models.CategoryPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *CourseController) EidtCategory() {
	if c.Ctx.Request.Method == "POST" {
		c.saveCategory()
	}
	c.setTpl("coursecategory/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "coursecategory/edit_footerjs.html"
	c.LayoutSections["headcssjs"] = "coursecategory/edit_headcssjs.html"
}

func (c *CourseController) saveCategory() {
	category := models.CourseCategory{}
	if err := c.ParseForm(&category); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", category.Id)
	}
	if category.Id > 0 {
		o := orm.NewOrm()
		if _, err := o.Update(&category); err != nil {
			c.jsonResult(enums.JRCodeFailed, "数据获取失败", category.Id)
		}
		c.jsonResult(enums.JRCodeSucc, "课程分类修改成功", category.Id)
	} else {
		o := orm.NewOrm()
		if _, err := o.Insert(&category); err != nil {
			c.jsonResult(enums.JRCodeFailed, "添加分类失败", category.Id)
		}
		c.jsonResult(enums.JRCodeSucc, "添加分类成功", category.Id)
	}
}

func (c *CourseController) DeleteCategory() {
	strs := c.GetString("ids")
	ids := make([]int, len(strs))
	for _, val := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(val); err == nil {
			ids = append(ids, id)
		}
	}
	if num, err := models.CategoryBatchDelete(ids); err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除失败", err)
	} else {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	}
}

/*
func (c *CourseController) GetClassify() {
	list := models.CourseAll(c.curUser.Id, 100)
	c.jsonResult(enums.JRCodeSucc, "", list)
}

func (c *CourseController) UpdateSeq() {
	// if c.Ctx.Request.Method == "POST" {
	// 	return
	// }
	// Id, _ := c.GetInt(":id", 0)
	// if Id == 0 {
	// 	c.jsonResult(enums.JRCodeFailed, "选择的数据无效", 0)
	// }

	// if m, err := models.AscienceOne(Id); err == nil {
	// 	c.Data[m] = m
	// 	// c.jsonResult(enums.JRCodeSucc, "", 0)
	// } else {
	// 	// c.jsonResult(enums.JRCodeFailed, "", 0)
	// }
	// c.setTpl("course/edit.html", "shared/layout_pullbox.html")
}



func getRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//不能判断一定是，可以判断一定不是。判断方式，base64只包含特定字符;解码再转码，查验是否相等。目前貌似没有能一定判断是的方法，有的话请指正，感谢。
func judgeBase64(str string) bool {
	pattern := "^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{4}|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)$"
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		return false
	}
	if !(len(str)%4 == 0 && matched) {
		return false
	}
	unCodeStr, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return false
	}
	tranStr := base64.StdEncoding.EncodeToString(unCodeStr)
	//return str==base64.StdEncoding.EncodeToString(unCodeStr)
	if str == tranStr {
		return true
	}
	return false
}
*/
