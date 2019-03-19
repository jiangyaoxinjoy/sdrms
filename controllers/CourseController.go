package controllers

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
	Tag, _ := c.GetInt(":tag", 0)
	c.Data["tag"] = Tag
	c.setTpl()
}
