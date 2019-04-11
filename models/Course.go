package models

import (
	"fmt"
	"time"

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
	CreatedTime    time.Time       `orm:"auto_now_add;type(datetime)"`
	ModifyTime     time.Time       `orm:"auto_now;type(datetime)"`
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

func TagsPageList(params *BaseQueryParam) ([]*Tag, int64) {
	o := orm.NewOrm()
	data := make([]*Tag, 0)
	sortorder := "Id"
	total, _ := o.QueryTable(TagTBName()).OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	// for _, val := range data {
	// 	if err := o.Read(val); err == nil {
	// 		o.LoadRelated(val, "Courses")
	// 		fmt.Println(val)
	// 	}
	// }
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

// //获取一个课程的信息
// func OneCourseDetail(id int) (*Course, *CourseContent) {
// 	m := &Course{Id: id}
// 	o := orm.NewOrm()
// 	o.Read(m, "Id")

// 	c := &CourseContent{Id: m.CourseContent.Id}
// 	o.Read(c, "Id")
// 	return m, c
// }

/*
type CourseResource struct {
	Title      string
	UserName   string
	CourseName string
	Id         int
	UploaderId int
}

type PostCourseResource struct {
	Title      string
	UserName   string
	CourseName string
	Modify     int
	Id         int
}

type PostMsg struct {
	List  []*PostCourseResource
	Pages *Page
}

type Page struct {
	PageNo     int    //当前页
	PageSize   int    //每页多少数据
	TotalPage  int    //总共多少页
	TotalCount int    //总共多少条数据
	FirstPage  int    //第一页
	LastPage   int    //最后一页
	Url        string //链接
	PageList   []int
	CourseId   int
}
*/

/*
func AscienceTreeGrid(courseid, cid, curpage int) *PostMsg {
	var list []*CourseResource
	var count int
	o := orm.NewOrm()
	page := new(Page)
	page.PageNo = curpage
	if courseid == 0 {
		o.Raw("SELECT COUNT(*) FROM rms_course_detail").QueryRow(&count)
		fmt.Println("*****************************")
		fmt.Println(count)
		page.TotalCount = count
		page.PageSize = 5
		page.getPageCount()
		fmt.Println(page)
		sql := fmt.Sprintf(`SELECT T0.title,T1.user_name,T2.course_name,T0.id,T0.uploader_id
		FROM %s AS T0
		INNER JOIN %s AS T1 ON T0.uploader_id = T1.id
		INNER JOIN %s AS T2 ON T2.id = T0.course_id
		LIMIT ? ,?
		 `, AscienceTBName(), BackendUserTBName(), CourseTBName())
		o.Raw(sql, (curpage-1)*5+1, curpage*5).QueryRows(&list)
	} else {
		o.Raw("SELECT COUNT(*) FROM rms_course_detail WHERE course_id = ?", courseid).QueryRow(&count)
		fmt.Println("*****************************")
		fmt.Println(count)
		page.PageSize = 5
		page.TotalCount = count
		page.getPageCount()

		sql := fmt.Sprintf(`SELECT DISTINCT T0.title,T1.user_name,T2.course_name,T0.id,T0.uploader_id
		FROM %s AS T0
		INNER JOIN %s AS T1 ON T0.uploader_id = T1.id
		INNER JOIN %s AS T2 ON T2.id = T0.course_id
		WHERE T0.course_id = ?
		LIMIT ? ,?
		 `, AscienceTBName(), BackendUserTBName(), CourseTBName())
		total, _ := o.Raw(sql, courseid, (curpage-1)*5+1, curpage*5).QueryRows(&list)
		fmt.Println("*****************************")
		fmt.Println(total)

	}
	var strlist []*PostCourseResource
	for _, val := range list {
		var modify int
		if val.UploaderId == cid {
			modify = 1
		} else {
			modify = 0
		}
		ins := new(PostCourseResource)
		ins.Title = val.Title
		ins.CourseName = val.CourseName
		ins.UserName = val.UserName
		ins.Modify = modify
		ins.Id = val.Id
		strlist = append(strlist, ins)
	}
	page.CourseId = courseid
	postmsg := new(PostMsg)
	postmsg.List = strlist
	postmsg.Pages = page
	return postmsg
}

func AscienceDeleteOne(id, curuserid int) int64 {
	o := orm.NewOrm()
	sql := fmt.Sprintf(`DELETE FROM %s WHERE id = ? AND uploader_id = ?`, AscienceTBName())
	res, err := o.Raw(sql, id, curuserid).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		fmt.Println("mysql row affected nums: ", num)
		return num
	} else {
		return 0
	}

}

func AscienceOne(id int) *Ascience {
	m := new(Ascience)
	o := orm.NewOrm()
	sql := fmt.Sprintf(`SELECT * FROM %s
		WHERE id = ?
		 `, AscienceTBName())
	o.Raw(sql, id).QueryRow(&m)
	return m
}

func AscienceUpdate(m *Ascience, uploadId int) (string, error) {
	o := orm.NewOrm()
	logs.Info("======+++++++++")
	// logs.Info(m)
	if m.Id == 0 {
		m.UploaderId = uploadId
		if _, err := o.Insert(m); err != nil {
			return "添加失败", err
		}
		return "添加成功", nil

	} else {
		if _, err := o.Update(m); err != nil {
			return "编辑失败", err
		}
		return "编辑成功", nil
	}

}

func (this *Page) getPageCount() {
	var tp float32 = float32(this.TotalCount) / float32(this.PageSize)
	if tp < 1 {
		this.TotalPage = 1
	}
	var tpint float32 = float32(int(tp))

	if tp > tpint {
		tpint += 1
	}
	this.TotalPage = int(tpint)
	this.LastPage = int(tpint)
	this.FirstPage = 1
	var i int = 1
	for ; ; i++ {
		this.PageList = append(this.PageList, i)
		if i >= this.TotalPage {
			break
		}
	}
}
*/
