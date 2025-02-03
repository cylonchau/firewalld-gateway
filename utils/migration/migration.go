package migration

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/cylonchau/firewalld-gateway/config"
	"github.com/cylonchau/firewalld-gateway/server/app/router"
	"github.com/cylonchau/firewalld-gateway/utils/model"
)

func Migration(driver string) error {
	if driver == "" {
		return errors.New("Unkown database driver")
	}
	var (
		dbInterface   *gorm.DB
		enconterError error
	)
	switch driver {
	case "mysql":
		if dbInterface, enconterError = MySQL(); enconterError == nil {
			// 执行迁移
			if err := autoMigrate(dbInterface); err != nil {
				return err
			}
		}
		return enconterError
	case "sqlite":
		if dbInterface, enconterError = SQLite(); enconterError == nil {
			if err := autoMigrate(dbInterface); err != nil {
				return err
			}
		}
		return nil
	}
	return enconterError
}

func initialData(db *gorm.DB) {
	// inital roles
	var (
		user_w_router_ids, user_r_router_ids         []model.Router // /security/user
		token_w_router_ids, token_r_router_ids       []model.Router // token
		host_w_router_ids, host_r_router_ids         []model.Router // /fw/host/
		tag_w_router_ids, tag_r_router_ids           []model.Router // /fw/tag/
		auth_w_router_ids, auth_r_router_ids         []model.Router // /security/auth
		service_r_router_ids, service_w_router_ids   []model.Router // service
		nat_r_router_ids, nat_w_router_ids           []model.Router // NAT
		port_r_router_ids, port_w_router_ids         []model.Router // port
		rich_r_router_ids, rich_w_router_ids         []model.Router // rich
		setting_w_router_ids, setting_r_router_ids   []model.Router // setting
		template_w_router_ids, template_r_router_ids []model.Router
		audit_r_router_ids                           []model.Router // audit
	)

	db.Select([]string{"id"}).Where("path LIKE '%/security/users%' AND method != 'GET'").Find(&user_w_router_ids)                   // user
	db.Select([]string{"id"}).Where("path LIKE '%/security/users%' AND method = 'GET'").Find(&user_r_router_ids)                    // user
	db.Select([]string{"id"}).Where("path LIKE '%/security/tokens%' AND method != 'GET'").Find(&token_w_router_ids)                 // token
	db.Select([]string{"id"}).Where("path LIKE '%/security/tokens%' AND method = 'GET'").Find(&token_r_router_ids)                  // token
	db.Select([]string{"id"}).Where("path LIKE '%/fw/host%' AND method != 'GET'").Find(&host_w_router_ids)                          // host
	db.Select([]string{"id"}).Where("path LIKE '%/fw/host%' AND method = 'GET'").Find(&host_r_router_ids)                           // host
	db.Select([]string{"id"}).Where("path LIKE '%/fw/tag%' AND method != 'GET'").Find(&tag_w_router_ids)                            // tag
	db.Select([]string{"id"}).Where("path LIKE '%/fw/tag%' AND method = 'GET'").Find(&tag_r_router_ids)                             // tag
	db.Select([]string{"id"}).Where("path LIKE '%/security/auth%' AND method != 'GET'").Find(&auth_w_router_ids)                    // auth
	db.Select([]string{"id"}).Where("path LIKE '%/security/auth%' AND method = 'GET'").Find(&auth_r_router_ids)                     // auth
	db.Select([]string{"id"}).Where("path LIKE '%/service%' AND method != 'GET'").Find(&service_w_router_ids)                       // fw service
	db.Select([]string{"id"}).Where("path LIKE '%/service%' AND method = 'GET'").Find(&service_r_router_ids)                        // fw service
	db.Select([]string{"id"}).Where("path LIKE '%/nat%' AND path LIKE '%/masquerade%' AND method != 'GET'").Find(&nat_w_router_ids) // fw NAT
	db.Select([]string{"id"}).Where("path LIKE '%/nat%' AND path LIKE '%/masquerade%' AND method = 'GET'").Find(&nat_r_router_ids)  // fw NAT
	db.Select([]string{"id"}).Where("path LIKE '%/ports%' AND method != 'GET'").Find(&port_w_router_ids)                            // fw port
	db.Select([]string{"id"}).Where("path LIKE '%/ports%' AND method = 'GET'").Find(&port_r_router_ids)                             // fw port
	db.Select([]string{"id"}).Where("path LIKE '%/rich%' AND method != 'GET'").Find(&rich_w_router_ids)                             // fw rich
	db.Select([]string{"id"}).Where("path LIKE '%/rich%' AND method = 'GET'").Find(&rich_r_router_ids)                              // fw rich
	db.Select([]string{"id"}).Where("path LIKE '%/rich%' AND method != 'GET'").Find(&setting_w_router_ids)                          // fw setting
	db.Select([]string{"id"}).Where("path LIKE '%/rich%' AND method = 'GET'").Find(&setting_r_router_ids)                           // fw setting
	db.Select([]string{"id"}).Where("path LIKE '%/audit%' and method = 'GET'").Find(&audit_r_router_ids)                            // fw audit
	db.Select([]string{"id"}).Where("path LIKE '%/rich%' AND method = 'GET'").Find(&setting_r_router_ids)                           // fw setting
	db.Select([]string{"id"}).Where("path LIKE '%/audit%' and method = 'GET'").Find(&audit_r_router_ids)                            // fw audit
	db.Create(&model.Role{Name: "user_editer", Routers: user_w_router_ids})
	db.Create(&model.Role{Name: "user_viewer", Routers: user_r_router_ids})
	db.Create(&model.Role{Name: "token_viewer", Routers: token_r_router_ids})
	db.Create(&model.Role{Name: "token_editer", Routers: token_w_router_ids})
	db.Create(&model.Role{Name: "host_editer", Routers: host_w_router_ids})
	db.Create(&model.Role{Name: "host_viewer", Routers: host_r_router_ids})
	db.Create(&model.Role{Name: "tag_editer", Routers: tag_w_router_ids})
	db.Create(&model.Role{Name: "tag_viewer", Routers: tag_r_router_ids})
	db.Create(&model.Role{Name: "auth_viewer", Routers: auth_r_router_ids})
	db.Create(&model.Role{Name: "auth_editer", Routers: auth_w_router_ids})
	db.Create(&model.Role{Name: "service_viewer", Routers: service_r_router_ids})
	db.Create(&model.Role{Name: "service_editer", Routers: service_w_router_ids})
	db.Create(&model.Role{Name: "nat_viewer", Routers: nat_r_router_ids})
	db.Create(&model.Role{Name: "nat_editer", Routers: nat_w_router_ids})
	db.Create(&model.Role{Name: "port_viewer", Routers: port_r_router_ids})
	db.Create(&model.Role{Name: "port_editer", Routers: port_w_router_ids})
	db.Create(&model.Role{Name: "rich_viewer", Routers: rich_r_router_ids})
	db.Create(&model.Role{Name: "rich_editer", Routers: rich_w_router_ids})
	db.Create(&model.Role{Name: "setting_viewer", Routers: setting_r_router_ids})
	db.Create(&model.Role{Name: "setting_editer", Routers: setting_w_router_ids})
	db.Create(&model.Role{Name: "template_viewer", Routers: template_r_router_ids})
	db.Create(&model.Role{Name: "template_editer", Routers: template_w_router_ids})
	db.Create(&model.Role{Name: "audit_viewer", Routers: audit_r_router_ids})

	db.Create(&model.User{Username: "admin", Password: model.EncryptPassword("admin")})
}

func autoMigrate(dbInterface *gorm.DB) (enconterError error) {
	dbInterface.Migrator().CurrentDatabase()
	if !dbInterface.Migrator().HasTable(&model.User{}) {
		if enconterError = dbInterface.Migrator().AutoMigrate(&model.User{}); enconterError != nil {
			return enconterError
		}
	}
	if !dbInterface.Migrator().HasTable(&model.Tag{}) {
		if enconterError = dbInterface.Migrator().AutoMigrate(&model.Tag{}); enconterError != nil {
			return enconterError
		}
	}
	if !dbInterface.Migrator().HasTable(&model.Host{}) {
		if enconterError = dbInterface.Migrator().AutoMigrate(&model.Host{}); enconterError != nil {
			return enconterError
		}
	}
	if !dbInterface.Migrator().HasTable(&model.Template{}) {
		if enconterError = dbInterface.Migrator().AutoMigrate(&model.Template{}); enconterError != nil {
			return enconterError
		}
	}
	if !dbInterface.Migrator().HasTable(&model.Port{}) {
		if enconterError = dbInterface.Migrator().AutoMigrate(&model.Port{}); enconterError != nil {
			return enconterError
		}
	}
	if !dbInterface.Migrator().HasTable(&model.Rich{}) {
		if enconterError = dbInterface.Migrator().AutoMigrate(&model.Rich{}); enconterError != nil {
			return enconterError
		}

	}
	if !dbInterface.Migrator().HasTable(&model.Token{}) {
		if enconterError = dbInterface.Migrator().AutoMigrate(&model.Token{}); enconterError != nil {
			return enconterError
		}
	}
	if !dbInterface.Migrator().HasTable(&model.Audit{}) {
		if enconterError = dbInterface.Migrator().AutoMigrate(&model.Audit{}); enconterError != nil {
			return enconterError
		}
	}

	if !dbInterface.Migrator().HasTable(&model.Role{}) || !dbInterface.Migrator().HasTable(&model.Router{}) {
		if enconterError = dbInterface.AutoMigrate(&model.Role{}, &model.Router{}); enconterError != nil {
			return enconterError
		}
		http := gin.New()
		router.RegisteredRouter(http)

		for _, route := range http.Routes() {
			result := dbInterface.Create(&model.Router{
				Path:   route.Path,
				Method: route.Method,
			})
			fmt.Println(route.Path)
			if enconterError = result.Error; enconterError != nil {
				return enconterError
			}
		}
		initialData(dbInterface)
	}
	return nil
}

func SQLite() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(config.CONFIG.SQLite.File+".db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
}

func MySQL() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.CONFIG.MySQL.User,
		config.CONFIG.MySQL.Password,
		config.CONFIG.MySQL.IP,
		config.CONFIG.MySQL.Port,
		config.CONFIG.MySQL.Database,
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
}
