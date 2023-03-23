package batch_processor

import (
	"reflect"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalld-gateway/config"
)

var P *Processor
var T = time.Second * 5

type Processor struct {
	stopCh        chan interface{}
	addCh         chan interface{}
	listenersLock sync.RWMutex
	wg            wait.Group
	queue         workqueue.RateLimitingInterface
}

func NewProcessor() *Processor {
	if !reflect.DeepEqual(P, nil) {
		return &Processor{
			stopCh: make(chan interface{}),
			addCh:  make(chan interface{}),
			queue:  workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		}
	}
	return P
}

func (p *Processor) Run() {
	defer func() {
		p.queue.ShutDown()
	}()
	p.wg.Start(p.pop)
	p.wg.Wait()
}

func (p *Processor) Add(notification string, event interface{}) {
	StoreAdd(notification, event)

	p.queue.Add(notification)
}

func (p *Processor) AddAfter(notification string, t time.Duration, event interface{}) {
	if v, ok := event.(interface{}); ok && v != nil {
		StoreAdd(notification, event)
	}
	p.queue.AddAfter(notification, t)
}

func (p *Processor) pop() {

	klog.V(5).Infof("Async event processor started, waitting task...")

	for {
		var event Event
		select {
		case <-p.stopCh:
			klog.V(5).Infof("Async evnet process exit.")
			return
		default:
			notificationKey, quit := p.queue.Get()
			if quit {
				return
			}

			var encouterError error
			key := notificationKey.(string)
			eventInterface, enconterBool := Store[key]
			if enconterBool {
				event, enconterBool = eventInterface.(Event)
				if enconterBool {
					klog.V(5).Infof("Recived mission %s", event.TaskName)
					encouterError = event.processEvent()
					if encouterError != nil {
						if event.errNum <= config.CONFIG.Mission_Retry_Number {
							event.errNum++
							Store[key] = event
							retryTime := time.Duration(event.errNum+1) * T
							p.queue.Forget(key)
							p.AddAfter(key, retryTime, nil)
							klog.Warningf("Event processing failed, will retry on %v second after.", retryTime)
						} else {
							p.queue.Forget(key)
							StoreDel(key)
							klog.Warningf("Task %s exceed MRN value: %v.", event.TaskName, encouterError)
						}
					} else {
						p.queue.Forget(key)
						StoreDel(key)
					}
				}
			}
			p.queue.Done(key)
			if encouterError != nil || !enconterBool {
				klog.Errorf("Event failed: %v", encouterError)
			}
		}
	}
}
