package router

import (
	"embed"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
)

// ping godoc
// @Summary Return process status.
// @Description Return process status.
// @Tags Health
// @Accept  json
// @Produce json
// @Success 200 {object} interface{}
// @Router /ping [get]
func ping(c *gin.Context) {
	query.SuccessResponse(c, query.OK, "pong")
}

//go:embed dist/*
var distHandle embed.FS

//func dist() http.FileSystem {
//	//fsys, err := fs.Sub(distHandle, "dist")
//	//if err != nil {
//	//	panic(err)
//	//}
//	return http.FS(fsys)
//}
