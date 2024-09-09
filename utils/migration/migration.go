package migration

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

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
	case "sqlite":
		if dbInterface, enconterError = SQLite(); enconterError == nil {
			dbInterface.Migrator().CurrentDatabase()
			if !dbInterface.Migrator().HasTable(&model.User{}) {
				dbInterface.Migrator().AutoMigrate(&model.User{})
			}
			if !dbInterface.Migrator().HasTable(&model.Tag{}) {
				dbInterface.Migrator().AutoMigrate(&model.Tag{})
			}
			if !dbInterface.Migrator().HasTable(&model.Host{}) {
				dbInterface.Migrator().AutoMigrate(&model.Host{})
			}
			if !dbInterface.Migrator().HasTable(&model.Template{}) {
				dbInterface.Migrator().AutoMigrate(&model.Template{})
			}
			if !dbInterface.Migrator().HasTable(&model.Port{}) {
				dbInterface.Migrator().AutoMigrate(&model.Port{})
			}
			if !dbInterface.Migrator().HasTable(&model.Rich{}) {
				dbInterface.Migrator().AutoMigrate(&model.Rich{})
			}
			if !dbInterface.Migrator().HasTable(&model.Token{}) {
				dbInterface.Migrator().AutoMigrate(&model.Token{})
			}
			if !dbInterface.Migrator().HasTable(&model.Audit{}) {
				dbInterface.Migrator().AutoMigrate(&model.Audit{})
			}

			if !dbInterface.Migrator().HasTable(&model.Role{}) || !dbInterface.Migrator().HasTable(&model.Router{}) {
				dbInterface.AutoMigrate(&model.Role{}, &model.Router{})
				http := gin.New()
				router.RegisteredRouter(http)
				for _, route := range http.Routes() {
					result := dbInterface.Create(&model.Router{
						Path:   route.Path,
						Method: route.Method,
					})
					if enconterError = result.Error; enconterError != nil {
						return enconterError
					}
				}
				initialData(dbInterface)
			}
		}
		return nil
	}
	return enconterError
}

func initialData(db *gorm.DB) {
	db.Create(&model.User{
		Username: "admin",
		Password: model.EncryptPassword("111"),
	})

	// inital roles
	var (
		user_w_router_ids, user_r_router_ids             []model.Router
		host_w_router_ids, host_r_router_ids             []model.Router
		tag_w_router_ids, tag_r_router_ids               []model.Router
		template_w_router_ids, template_r_router_ids     []model.Router
		masquerade_w_router_ids, masquerade_r_router_ids []model.Router
		setting_w_router_ids, setting_r_router_ids       []model.Router
		auth_w_router_ids, auth_r_router_ids             []model.Router
		audit_r_router_ids                               []model.Router
	)
	db.Select([]string{"id"}).Where("path like '%user%' and method != 'GET'").Find(&user_w_router_ids)
	db.Select([]string{"id"}).Where("path like '%user%' and method = 'GET'").Find(&user_r_router_ids)
	db.Select([]string{"id"}).Where("path like '%/fw/host%' and method != 'GET'").Find(&host_w_router_ids)
	db.Select([]string{"id"}).Where("path like '%/fw/host%' and method = 'GET'").Find(&host_r_router_ids)
	db.Select([]string{"id"}).Where("path like '%/fw/tag%' and method != 'GET'").Find(&tag_w_router_ids)
	db.Select([]string{"id"}).Where("path like '%/fw/tag%' and method = 'GET'").Find(&tag_r_router_ids)
	db.Select([]string{"id"}).Where("path like '%/fw/host%' and method != 'GET'").Find(&host_w_router_ids)
	db.Select([]string{"id"}).Where("path like '%/fw/host%' and method = 'GET'").Find(&host_r_router_ids)
	db.Select([]string{"id"}).Where("path like '%/masquerade%' and method != 'GET'").Find(&masquerade_w_router_ids)
	db.Select([]string{"id"}).Where("path like '%/masquerade%' and method = 'GET'").Find(&masquerade_r_router_ids)
	db.Select([]string{"id"}).Where("path like '%/audit%' and method = 'GET'").Find(&audit_r_router_ids)
	db.Select([]string{"id"}).Where("path like '%/auth%' and method != 'GET'").Find(&auth_w_router_ids)
	db.Select([]string{"id"}).Where("path like '%/auth%' and method = 'GET'").Find(&auth_r_router_ids)

	db.Create(&model.Role{Name: "user_editer", Routers: user_w_router_ids})
	db.Create(&model.Role{Name: "user_viewer", Routers: user_r_router_ids})

	db.Create(&model.Role{Name: "host_editer", Routers: host_w_router_ids})
	db.Create(&model.Role{Name: "host_viewer", Routers: host_r_router_ids})

	db.Create(&model.Role{Name: "tag_editer", Routers: tag_w_router_ids})
	db.Create(&model.Role{Name: "tag_viewer", Routers: tag_r_router_ids})

	db.Create(&model.Role{Name: "template_editer", Routers: template_w_router_ids})
	db.Create(&model.Role{Name: "template_viewer", Routers: template_r_router_ids})

	db.Create(&model.Role{Name: "masquerade_editer", Routers: masquerade_w_router_ids})
	db.Create(&model.Role{Name: "masquerade_viewer", Routers: masquerade_r_router_ids})

	db.Create(&model.Role{Name: "setting_editer", Routers: setting_w_router_ids})
	db.Create(&model.Role{Name: "setting_viewer", Routers: setting_r_router_ids})

	db.Create(&model.Role{Name: "audit_viewer", Routers: audit_r_router_ids})

	db.Create(&model.User{Username: "admin", Password: model.EncryptPassword("admin")})
}

func SQLite() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(config.CONFIG.SQLite.File+".db"), &gorm.Config{})
}

func MySQL() (*gorm.DB, error) {
	return nil, nil
}
