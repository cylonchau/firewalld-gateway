package batch_processor

import (
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

type Processor struct {
	nextCh        chan Event
	addCh         chan Event
	listenersLock sync.RWMutex
	wg            sync.WaitGroup
}

func (p *Processor) Start() {
	p.wg.Add(2)
	go func() {
		defer p.wg.Done()
		p.run()
	}()
	go func() {
		defer p.wg.Done()
		p.pop()
	}()
}

func (p *Processor) add(notification Event) {
	p.addCh <- notification
}

func (p *Processor) pop() {
	defer close(p.nextCh) // Tell .run() to stop

	//var nextCh chan<- interface{}
	//var notification interface{}
	//for {
	//	select {
	//	case nextCh <- notification:
	//		// Notification dispatched
	//		var ok bool
	//		notification, ok = p.pendingNotifications.ReadOne()
	//		if !ok { // Nothing to pop
	//			nextCh = nil // Disable this select case
	//		}
	//	case notificationToAdd, ok := <-p.addCh:
	//		if !ok {
	//			return
	//		}
	//		if notification == nil { // No notification to pop (and pendingNotifications is empty)
	//			// Optimize the case - skip adding to pendingNotifications
	//			notification = notificationToAdd
	//			nextCh = p.nextCh
	//		} else { // There is already a notification waiting to be dispatched
	//			p.pendingNotifications.WriteOne(notificationToAdd)
	//		}
	//	}
	//}
}

func (p *Processor) run() {
	// this call blocks until the channel is closed.  When a panic happens during the notification
	// we will catch it, **the offending item will be skipped!**, and after a short delay (one second)
	// the next notification will be attempted.  This is usually better than the alternative of never
	// delivering again.
	stopCh := make(chan struct{})
	wait.Until(func() {
		for next := range p.nextCh {
			next = next
		}
		// the only way to get here is if the p.nextCh is empty and closed
		close(stopCh)
	}, 1*time.Second, stopCh)
}

func (p *Processor) distribute(obj interface{}) {
	p.listenersLock.RLock()
	defer p.listenersLock.RUnlock()

}
