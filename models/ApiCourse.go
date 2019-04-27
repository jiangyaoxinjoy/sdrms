package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/dgrijalva/jwt-go"
)

func (a *Comment) TableName() string {
	return CommentTBName()
}

type jwtCustomClaims struct {
	jwt.StandardClaims

	// 追加自己需要的信息
	Uid   uint `json:"uid"`
	Admin bool `json:"admin"`
}

type ApiCourse struct {
	Id          int      `json:"aid"`
	Title       string   `json:"title"`
	Status      int      `json:"status"`
	Identify    string   `json:"identify"`
	CreatedTime string   `json:"date"`
	Author      string   `json:"author"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
	CommentNum  int      `json:"comment_n"`
}

type ApiCourseQueryParam struct {
	BaseQueryParam
	TagId int `json:"tagid"`
	Page  int `json:"page"`
}

type ApiCourseCategory struct {
	Id   int
	Name string
}

type ApiOneCourse struct {
	Id         int                `json:"aid"`
	Title      string             `json:"title"`
	Content    string             `json:"content"`
	Createtime string             `json:"date"`
	Username   string             `json:"name"`
	Category   *ApiCourseCategory `json:"tags"`
}

type ApiCategory struct {
	Id          int
	Name        string
	TotalCourse int64
}

type User struct {
	Username string `json:"name"`
	Userpwd  string `json:"password"`
}

type ApiUser struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type ApiTags struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type Comment struct {
	Id          int `orm:"pk;column(id);`
	Content     string
	Course      *Course `orm:"rel(fk)"` //设置一对多关系
	NickName    string
	CreatedTime string
	Like        int
	Addresss    string
	ImgName     string
}

type ApiComment struct {
	Id          int    `json:"id"`
	Content     string `json:"content"`
	CreatedTime string `json:"date"`
	Like        int    `json:"like"`
	ArticleId   string `json:"articleId"`
	Addresss    string `json:"address"`
	ImgName     string `json:"imgName"`
	NickName    string `json:"name"`
}

func ApiCoursePageList(params *ApiCourseQueryParam) ([]*Course, error) {
	data := make([]*Course, 0)
	o := orm.NewOrm()
	query := o.QueryTable(CourseTBName()).Filter("Status", 1)
	//默认排序 id降序
	sortorder := "-Id"

	if params.TagId > 0 {
		fmt.Println(params)
		query = query.Filter("CourseCategory__Id", params.TagId)
	}
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).RelatedSel("BackendUser", "CourseContent").All(&data)
	for _, val := range data {
		o.LoadRelated(val, "Tags")
		o.LoadRelated(val, "Comments")
	}
	return data, nil
}

func ApiGetCourseById(id int) (*Course, error) {
	o := orm.NewOrm()
	m := &Course{Id: id}
	o.QueryTable(CourseTBName()).Filter("Id", id).RelatedSel("CourseCategory", "BackendUser").One(m)
	content := &CourseContent{}
	o.QueryTable(CourseContentTBName()).Filter("Course__id", id).One(content)
	m.CourseContent = content
	comments := make([]*Comment, 0)
	o.QueryTable(CommentTBName()).Filter("Course__id", id).All(&comments)
	m.Comments = comments
	o.LoadRelated(m, "Tags")
	return m, nil
}

func ApiGetCategorys() (*[]ApiCategory, error) {
	o := orm.NewOrm()
	categorys := make([]*CourseCategory, 0)
	num, err := o.QueryTable(CourseCategoryTBName()).Filter("Status", 1).RelatedSel().All(&categorys)
	if err != nil {
		return nil, err
	}
	data := make([]ApiCategory, num)
	for i, val := range categorys {
		fmt.Println(val)
		if n, err := o.QueryTable(CourseTBName()).Filter("CourseCategory__Id", val.Id).Count(); err == nil {
			data[i].Id = val.Id
			data[i].Name = val.Name
			data[i].TotalCourse = n
		}
	}
	return &data, nil
}

func ApiGetOneCategoryById(id int) (*CourseCategory, error) {
	courseCategory := new(CourseCategory)
	courseCategory.Id = id
	o := orm.NewOrm()
	o.Read(courseCategory)
	return courseCategory, nil
}

/**
 * 生成 token
 * SecretKey 是一个 const 常量
 */
// func CreateToken(SecretKey []byte, issuer string, Uid uint, isAdmin bool) (tokenString string, err error) {
// 	claims := &jwtCustomClaims{
// 		jwt.StandardClaims{
// 			ExpiresAt: int64(time.Now().Add(time.Hour * 72).Unix()),
// 			Issuer:    issuer,
// 		},
// 		Uid,
// 		isAdmin,
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err = token.SignedString(SecretKey)
// 	return
// }

func ApiCreateToken(key string, m map[string]interface{}) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)

	for index, val := range m {
		claims[index] = val
	}
	// fmt.Println(_map)
	token.Claims = claims
	tokenString, _ := token.SignedString([]byte(key))
	return tokenString
}
func ApiParseToken(tokenString string, key string) (interface{}, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		fmt.Println(err)
		return "", false
	}
}

func ApiGetCommentByIdentify(identify string) ([]*Comment, error) {
	o := orm.NewOrm()
	data := make([]*Comment, 0)
	o.QueryTable(CommentTBName()).Filter("Course__Identify", identify).All(&data)
	fmt.Println(data)
	return data, nil
}

func ApiAddComment(data *ApiComment) error {
	o := orm.NewOrm()
	comment := &Comment{}
	comment.Content = data.Content
	comment.Addresss = data.Addresss
	comment.CreatedTime = data.CreatedTime
	comment.NickName = data.NickName
	comment.ImgName = data.ImgName
	course := new(Course)
	o.QueryTable(CourseTBName()).Filter("Identify", data.ArticleId).One(course)
	comment.Course = course
	_, err := o.Insert(comment)
	return err
}
