package timewheel

import (
	"container/list"
	"sync"
	"time"
)

//标记一个任务的位置信息
type location struct {
	slot  int
	etask *list.Element
}

type TimeWheel struct {
	interval time.Duration //时间轮的刻度
	ticker   *time.Ticker  //每隔interval时间就往前推进
	slots    []*list.List  //每一个槽都是双向链表

	timer             map[string]*location //根据key快速拿到任务的地址
	currentPos        int                  //当前执行到几号槽位
	slotNum           int                  //槽的总数
	addTaskChannel    chan task
	removeTaskChannel chan string
	stopChannel       chan bool

	mu sync.RWMutex
}

type task struct {
	delay  time.Duration //多久后执行任务
	circle int           //执行的时机
	key    string        //任务的身份标识
	job    func()        //任务的具体执行函数
}

func New(interval time.Duration, slotNum int) *TimeWheel {
	if interval <= 0 || slotNum <= 0 {
		return nil
	}
	tw := &TimeWheel{
		interval:          interval,
		slots:             make([]*list.List, slotNum),
		timer:             make(map[string]*location),
		currentPos:        0,
		slotNum:           slotNum,
		addTaskChannel:    make(chan task),
		removeTaskChannel: make(chan string),
		stopChannel:       make(chan bool),
	}
	tw.initSlots()
	return tw
}
func (tw *TimeWheel) initSlots() {
	for i := 0; i < tw.slotNum; i++ {
		tw.slots[i] = list.New()
	}
}

func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.start()
}

func (tw *TimeWheel) Stop() {
	tw.stopChannel <- true
}

func (tw *TimeWheel) AddJob(delay time.Duration, key string, job func()) {
	if delay < 0 {
		return
	}
	tw.addTaskChannel <- task{delay: delay, key: key, job: job}
}

func (tw *TimeWheel) RemoveJob(key string) {
	if key == "" {
		return
	}
	tw.removeTaskChannel <- key
}

func (tw *TimeWheel) start() {
	//1.ticker有消息,taskhandler去执行
	//2.新增task
	//3.删除task
	for {
		select {
		case <-tw.ticker.C:
			tw.tickHandler()
		case task := <-tw.addTaskChannel:
			tw.addTask(&task)
		case key := <-tw.removeTaskChannel:
			tw.removeTask(key)
		case <-tw.stopChannel:
			tw.ticker.Stop()
			return
		}
	}
}

//前进一格，执行当前槽的全部任务
func (tw *TimeWheel) tickHandler() {
	tw.mu.Lock()
	l := tw.slots[tw.currentPos]
	if tw.currentPos == tw.slotNum-1 {
		tw.currentPos = 0
	} else {
		tw.currentPos++
	}
	tw.mu.Unlock()
	go tw.scanAndRunTask(l)
}

func (tw *TimeWheel) scanAndRunTask(l *list.List) {
	var tasksToRemove []string
	tw.mu.RLock()
	for e := l.Front(); e != nil; {
		task := e.Value.(*task)
		if task.circle > 0 {
			task.circle--
			e = e.Next()
			continue
		}
		//circle==0才执行
		//开一个goroutine去执行一个任务
		go func(job func()) {
			defer func() {
				if err := recover(); err != nil {
				}
			}()
			job()
		}(task.job)
		if task.key != "" {
			tasksToRemove = append(tasksToRemove, task.key)
		}
		//执行结束就删除
		next := e.Next()
		l.Remove(e)
		e = next
	}
	tw.mu.RUnlock()
	tw.mu.Lock()
	//删除所有执行结束的任务的位置信息
	for _, key := range tasksToRemove {
		delete(tw.timer, key)
	}
	tw.mu.Unlock()
}

func (tw *TimeWheel) addTask(task *task) {
	pos, circle := tw.getPositionAndCircle(task.delay)
	task.circle = circle
	//计算当前task的slot和circle
	tw.mu.Lock()
	defer tw.mu.Unlock()
	//如果相同的key，task已经存在就删除旧任务
	if task.key != "" {
		if _, ok := tw.timer[task.key]; ok {
			tw.removeTaskInternal(task.key)
		}
	}
	//在slot 的链表加上新节点
	//将位置信息保存在timer
	e := tw.slots[pos].PushBack(task)
	loc := &location{
		slot:  pos,
		etask: e,
	}
	tw.timer[task.key] = loc
}

func (tw *TimeWheel) getPositionAndCircle(d time.Duration) (pos int, circle int) {
	delaySeconds := int(d.Seconds())
	intervalSeconds := int(tw.interval.Seconds())
	//delaySeconds / intervalSeconds当前任务需要经过几个slot后执行
	//slots/slotNum得到circle
	circle = delaySeconds / intervalSeconds / tw.slotNum
	//当前slot+newslot就是当前多久后执行
	pos = (tw.currentPos + delaySeconds/intervalSeconds) % tw.slotNum
	return
}

func (tw *TimeWheel) removeTask(key string) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	tw.removeTaskInternal(key)
}

func (tw *TimeWheel) removeTaskInternal(key string) {
	//slot中的链表删除task节点
	//删除key对应的任务位置信息
	pos, ok := tw.timer[key]
	if !ok {
		return
	}
	l := tw.slots[pos.slot]
	l.Remove(pos.etask)
	delete(tw.timer, key)
}
