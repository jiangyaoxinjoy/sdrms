package routers

import (
	"sdrms/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//课程路由
	beego.Router("/course/index/", &controllers.CourseController{}, "*:Index")
	beego.Router("/course/datagrid", &controllers.CourseController{}, "Post:DataGrid")
	beego.Router("/course/tagsdatagrid", &controllers.CourseController{}, "Post:TagsDataGrid")
	beego.Router("/course/categorydatagrid", &controllers.CourseController{}, "Post:CategoryDataGrid")
	beego.Router("/course/delete", &controllers.CourseController{}, "Post:Delete")
	beego.Router("/course/edit/?:id", &controllers.CourseController{}, "Get,Post:Edit")
	beego.Router("/course/editormdPic", &controllers.CourseController{}, "Post:EditormdPic")
	beego.Router("/course/tags", &controllers.CourseController{}, "Get:Tags")
	beego.Router("/course/deletetag", &controllers.CourseController{}, "Post:DeleteTag")
	beego.Router("/course/edittag", &controllers.CourseController{}, "Post,Get:EidtTag")
	beego.Router("/course/category", &controllers.CourseController{}, "Get:Category")
	beego.Router("/course/editcategory", &controllers.CourseController{}, "Post,Get:EidtCategory")
	beego.Router("/course/deletecategory", &controllers.CourseController{}, "Post:DeleteCategory")

	//用户角色路由
	beego.Router("/role/index", &controllers.RoleController{}, "*:Index")
	beego.Router("/role/datagrid", &controllers.RoleController{}, "Get,Post:DataGrid")
	beego.Router("/role/edit/?:id", &controllers.RoleController{}, "Get,Post:Edit")
	beego.Router("/role/delete", &controllers.RoleController{}, "Post:Delete")
	beego.Router("/role/datalist", &controllers.RoleController{}, "Post:DataList")
	beego.Router("/role/allocate", &controllers.RoleController{}, "Post:Allocate")
	beego.Router("/role/updateseq", &controllers.RoleController{}, "Post:UpdateSeq")

	//资源路由
	beego.Router("/resource/index", &controllers.ResourceController{}, "*:Index")
	beego.Router("/resource/treegrid", &controllers.ResourceController{}, "POST:TreeGrid")
	beego.Router("/resource/edit/?:id", &controllers.ResourceController{}, "Get,Post:Edit")
	beego.Router("/resource/parent", &controllers.ResourceController{}, "Post:ParentTreeGrid")
	beego.Router("/resource/delete", &controllers.ResourceController{}, "Post:Delete")
	//快速修改顺序
	beego.Router("/resource/updateseq", &controllers.ResourceController{}, "Post:UpdateSeq")

	//通用选择面板
	beego.Router("/resource/select", &controllers.ResourceController{}, "Get:Select")
	//用户有权管理的菜单列表（包括区域）
	beego.Router("/resource/usermenutree", &controllers.ResourceController{}, "POST:UserMenuTree")
	beego.Router("/resource/checkurlfor", &controllers.ResourceController{}, "POST:CheckUrlFor")

	//后台用户路由
	beego.Router("/backenduser/index", &controllers.BackendUserController{}, "*:Index")
	beego.Router("/backenduser/datagrid", &controllers.BackendUserController{}, "POST:DataGrid")
	beego.Router("/backenduser/edit/?:id", &controllers.BackendUserController{}, "Get,Post:Edit")
	beego.Router("/backenduser/delete", &controllers.BackendUserController{}, "Post:Delete")
	//后台用户中心
	beego.Router("/usercenter/profile", &controllers.UserCenterController{}, "Get:Profile")
	beego.Router("/usercenter/basicinfosave", &controllers.UserCenterController{}, "Post:BasicInfoSave")
	beego.Router("/usercenter/uploadimage", &controllers.UserCenterController{}, "Post:UploadImage")
	beego.Router("/usercenter/passwordsave", &controllers.UserCenterController{}, "Post:PasswordSave")

	beego.Router("/home/index", &controllers.HomeController{}, "*:Index")
	beego.Router("/home/login", &controllers.HomeController{}, "*:Login")
	beego.Router("/home/dologin", &controllers.HomeController{}, "Post:DoLogin")
	beego.Router("/home/logout", &controllers.HomeController{}, "*:Logout")

	beego.Router("/home/404", &controllers.HomeController{}, "*:Page404")
	beego.Router("/home/error/?:error", &controllers.HomeController{}, "*:Error")

	beego.Router("/", &controllers.HomeController{}, "*:Index")

}
