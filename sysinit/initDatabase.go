package sysinit

import (
	_ "sdrms/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	//_ "github.com/mattn/go-sqlite3"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//初始化数据连接
func InitDatabase() {
	//读取配置文件，设置数据库参数
	//数据库类别
	dbType := beego.AppConfig.String("db_type")
	//连接名称
	dbAlias := beego.AppConfig.String(dbType + "::db_alias")
	//数据库名称
	dbName := beego.AppConfig.String(dbType + "::db_name")
	//数据库连接用户名
	dbUser := beego.AppConfig.String(dbType + "::db_user")
	//数据库连接用户名
	dbPwd := beego.AppConfig.String(dbType + "::db_pwd")
	//数据库IP（域名）
	dbHost := beego.AppConfig.String(dbType + "::db_host")
	//数据库端口
	dbPort := beego.AppConfig.String(dbType + "::db_port")
	switch dbType {
	case "sqlite3":
		orm.RegisterDataBase(dbAlias, dbType, dbName)
	case "mysql":
		dbCharset := beego.AppConfig.String(dbType + "::db_charset")
		// 参数1        数据库的别名，用来在 ORM 中切换数据库使用
		// 参数2        driverName
		// 参数3        对应的链接字符串
		// 参数4(可选)  设置最大空闲连接
		// 参数5(可选)  设置最大数据库连接 (go >= 1.2)
		orm.RegisterDataBase(dbAlias, dbType, dbUser+":"+dbPwd+"@tcp("+dbHost+":"+
			dbPort+")/"+dbName+"?charset="+dbCharset, 30)
	}
	//如果是开发模式，则显示命令信息
	isDev := (beego.AppConfig.String("runmode") == "dev")
	//自动建表
	orm.RunSyncdb("default", false, isDev)
	// 设置为 UTC 时间
	orm.DefaultTimeLoc = time.UTC
	if isDev {
		orm.Debug = isDev
	}
}
