package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// init 初始化
func init() {
	orm.RegisterModel(new(BackendUser), new(Tag), new(Resource), new(Role), new(RoleResourceRel), new(RoleBackendUserRel), new(Course), new(CourseCategory), new(CourseContent), new(Comment), new(Map), new(MapStatus))
}

// TableName 下面是统一的表名管理
func TableName(name string) string {
	prefix := beego.AppConfig.String("db_dt_prefix")
	return prefix + name
}

// BackendUserTBName 获取 BackendUser 对应的表名称
func BackendUserTBName() string {
	return TableName("backend_user")
}

// ResourceTBName 获取 Resource 对应的表名称
func ResourceTBName() string {
	return TableName("resource")
}

// RoleTBName 获取 Role 对应的表名称
func RoleTBName() string {
	return TableName("role")
}

// RoleResourceRelTBName 角色与资源多对多关系表
func RoleResourceRelTBName() string {
	return TableName("role_resource_rel")
}

// RoleBackendUserRelTBName 角色与用户多对多关系表
func RoleBackendUserRelTBName() string {
	return TableName("role_backenduser_rel")
}

//获取对应课程表
func CourseCategoryTBName() string {
	return TableName("course_category")
}

func CourseTBName() string {
	return TableName("course")
}

func CourseContentTBName() string {
	return TableName("course_content")
}

func TagTBName() string {
	return TableName("tag")
}

func CommentTBName() string {
	return TableName("comment")
}

func MapTBName() string {
	return TableName("map")
}

func MapStatusTBName() string {
	return TableName("mapstatus")
}
