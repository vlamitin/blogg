package notifier

type Notifier struct {
	subscribers   map[*Subscriber]bool
	subscribeCh   chan *Subscriber
	unsubscribeCh chan *Subscriber
	notifyCh      chan string
}

func NewNotifier(notifyCh chan string) *Notifier {
	return &Notifier{
		subscribers:   make(map[*Subscriber]bool),
		subscribeCh:   make(chan *Subscriber),
		unsubscribeCh: make(chan *Subscriber),
		notifyCh:      notifyCh,
	}
}

func (n *Notifier) Run() {
	for {
		select {
		case newSub := <-n.subscribeCh:
			n.subscribers[newSub] = true
		case sub := <-n.unsubscribeCh:
			if _, ok := n.subscribers[sub]; ok {
				delete(n.subscribers, sub)
				close(sub.sendCh)
			}
		case message := <-n.notifyCh:

			for sub := range n.subscribers {
				select {
				case sub.sendCh <- message:
				default:
					close(sub.sendCh)
					delete(n.subscribers, sub)
				}
			}
		}
	}
}
