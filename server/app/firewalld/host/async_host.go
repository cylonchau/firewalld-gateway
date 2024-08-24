package host

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praserx/ipconv"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalld-gateway/config"
	"github.com/cylonchau/firewalld-gateway/server/batch_processor"
	query2 "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/model"
)

type AsyncHost struct{}

func (h *AsyncHost) RegisterAsyncHostAPI(g *gin.RouterGroup) {
	g.POST("/async", h.createHost)
}

func (h *AsyncHost) createHost(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &query2.AsyncHostQuery{}
	enconterError = c.ShouldBindJSON(&query)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}

	ip, ipnet, err := net.ParseCIDR(query.IPRange)
	if err != nil {
		query2.APIResponse(c, err, nil)
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
		defer cancel()
		ch := make(chan bool)
		threadName := "thread_" + batch_processor.RandName()
		go func() {
			inc := func(ip net.IP) {
				for j := len(ip) - 1; j >= 0; j-- {
					ip[j]++
					if ip[j] > 0 {
						break
					}
				}
			}

			ping := func(ipaddr string) bool {
				conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", ipaddr, config.CONFIG.DbusPort), time.Second)
				if err != nil {
					klog.Errorf("[async host task] Add host failed: %v", err)
					return false
				}
				conn.Close()
				return true

			}

			for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
				if ip[len(ip)-1] == 0 || ip[len(ip)-1] == 255 { // 跳过 0 和 255 地址
					continue
				}
				if ping(ip.String()) {
					ipInt, _ := ipconv.IPv4ToInt(ip)
					host := &model.Host{
						IP:    ipInt,
						TagId: query.TagId,
					}
					if err := model.CreateHostWithHost(host); err != nil {
						klog.Errorf("[async host task] Insert to db failed: %v", err)
					}
				}
			}
			ch <- true
		}()

		select {
		case <-ch:
			// 如果异步操作完成，则打印完成信息
			klog.V(4).Infof("Thread %s finishd", threadName)
		case <-ctx.Done():
			klog.V(4).Infof("Thread %s exceeded timeout and was not able to finish.", threadName)
		}
	}()

	query2.SuccessResponse(c, query2.OK, nil)
}
