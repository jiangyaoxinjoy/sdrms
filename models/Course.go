package models

import (
	"fmt"
	//"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func (a *Course) TableName() string {
	return CourseTBName()
}

func (a *CourseContent) TableName() string {
	return CourseContentTBName()
}

//TableName 设置CourseCategory表名
func (a *CourseCategory) TableName() string {
	return CourseCategoryTBName()
}

func (a *Tag) TableName() string {
	return TagTBName()
}

// func (a *CourseTagRel) TableName() string {
// 	return CourseTagTBName()
// }

// 用于搜索的类
type CourseQueryParam struct {
	BaseQueryParam
	TitleLike  string
	Status     string
	CategoryId string
}

type Course struct {
	Id             int `orm:"pk;column(id);`
	Title          string
	Status         int
	BackendUser    *BackendUser    `orm:"rel(fk)"` //设置一对多关系
	Tags           []*Tag          `orm:"rel(m2m)"`
	CourseContent  *CourseContent  `orm:"rel(one)"` //设置一对一关系
	CourseCategory *CourseCategory `orm:"rel(fk)"`  //设置一对多关系
	CreatedTime    string
	ModifyTime     string
	Identify       string     `orm:"pk;`
	Comments       []*Comment `orm:"pk;reverse(many)"` // 设置一对多的反向关系
}

type CourseContent struct {
	Id      int
	Content string
	Course  *Course `orm:"reverse(one)"` // 设置一对一反向关系(可选)
}

type CourseCategory struct {
	Id      int `orm:"pk"`
	Name    string
	Status  int
	Courses []*Course `orm:"reverse(many)"` // 设置一对多的反向关系
	Cover   string
}

type Tag struct {
	Id      int
	Name    string
	Courses []*Course `orm:"reverse(many)"`
	Status  int
}

// type CourseTagRel struct {
// 	Id     int
// 	Tag    *Tag    `orm:"rel(fk)"`
// 	Course *Course `orm:"rel(fk)"`
// }

// type CourseMain struct {
// 	Id     int    `form:"id"`
// 	Title  string `form:"title"`
// 	Status int    `form:"status"`
// }

//course单条
func CourseOne(id int) (*Course, error) {
	o := orm.NewOrm()
	m := &Course{Id: id}
	var er error
	if er = o.QueryTable(CourseTBName()).Filter("Id", id).RelatedSel("CourseCategory", "BackendUser", "CourseContent").One(m); er == nil {
		o.Read(m.CourseContent)

		if _, err := o.LoadRelated(m, "Tags"); err != nil {
			return nil, err
		}
		return m, nil
	}
	return nil, er
}

//所有课程类型
func AllCategorys() []*CourseCategory {
	var list []*CourseCategory
	o := orm.NewOrm()
	if _, err := o.QueryTable(CourseCategoryTBName()).Filter("Status", 1).All(&list); err != nil {
		fmt.Println(err)
		return nil
	}
	return list
}

//所有的标签
func AllTags() []*Tag {
	var list []*Tag
	o := orm.NewOrm()
	if _, err := o.QueryTable(TagTBName()).Filter("Status", 1).All(&list); err != nil {
		fmt.Println(err)
		return nil
	}
	return list
}

//根据id获取标签
func GetTagIds(ids []int, tags *[]*Tag) error {
	query := orm.NewOrm().QueryTable(TagTBName())
	if _, err := query.Filter("id__in", ids).All(tags); err != nil {
		return err
	}
	return nil
}

// //批量删除
func CourseBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(CourseTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

//批量删除tag
func TagBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(TagTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

//批量删除category
func CategoryBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(CourseCategoryTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

func CategoryPageList(params *BaseQueryParam) ([]*CourseCategory, int64) {
	o := orm.NewOrm()
	data := make([]*CourseCategory, 0)
	sortorder := "Id"
	total, _ := o.QueryTable(CourseCategoryTBName()).OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

func TagsPageList(params *BaseQueryParam) ([]*Tag, int64) {
	o := orm.NewOrm()
	data := make([]*Tag, 0)
	sortorder := "Id"
	total, _ := o.QueryTable(TagTBName()).OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

//获取分页数据
func CoursePageList(params *CourseQueryParam) ([]*Course, int64) {
	o := orm.NewOrm()
	query := o.QueryTable(CourseTBName())
	data := make([]*Course, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	if len(params.Status) > 0 {
		query = query.Filter("Status", params.Status)
	}
	if len(params.CategoryId) > 0 {
		query = query.Filter("CourseCategory__Id", params.CategoryId)
	}
	query = query.Filter("title__contains", params.TitleLike)
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).RelatedSel("CourseCategory", "BackendUser").All(&data)
	return data, total
}
