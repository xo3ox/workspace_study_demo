package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/rpc"
	"sync"
	"time"
)

// 1 实现3节点选举
// 2 改造成分布式选举代码，加入RPC调用
// 3 完整代码 自主选举 日志复制

// 定义3节点的常量
const raftCount = 3

// 声明leader主节点对象
type Leader struct {
	Term     int // 任期
	LeaderId int // 领导者编号
}

// raft 声明
type Raft struct {
	mu              sync.Mutex // 锁
	me              int        // 节点编号
	currentTerm     int        // 当前任期
	votedFor        int        // 为哪个节点投票
	state           int        // 3个状态 0 follower 1 candidate 2 leader
	lastMessageTime int64      // 发送最后一条数据的时间
	currentLeader   int        // 设置当前的领导者
	message         chan bool  // 节点间发信息的通道
	electCh         chan bool  // 选举的通道
	heartBeat       chan bool  // 心跳信号的通道
	heartbeatRe     chan bool  // 返回心跳的通道
	timeout         int        // 超时时间
}

// 0 表示还没上任， -1 没有领导者编号
var leader = Leader{0, -1}

func main() {
	// 过程：有3个节点，最初都是follower
	// 若有candidate状态，进行投票和拉票
	// 会产生leader

	// 创建3个节点
	for i := 0; i < raftCount; i++ {
		// 创建3个raft节点
		Make(i)
	}

	// 加入服务端监听
	rpc.Register(new(Raft)) // 注册一个Raft服务
	// 服务器处理绑定到http协议上
	rpc.HandleHTTP()
	// 监听服务
	err := http.ListenAndServe(":8080",nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		;
	}
}

// 创建节点
func Make(me int) *Raft {
	rf := &Raft{}
	rf.me = me
	rf.votedFor = -1 // -1 刚创建时，谁都不投，此时节点刚创建
	rf.state = 0
	rf.timeout = 0
	rf.currentLeader = -1 // 还没有领导
	// 节点任期
	rf.setTerm(0)

	// 初始化通道
	rf.message = make(chan bool)     // 节点间发信息的通道
	rf.electCh = make(chan bool)     // 选举的通道
	rf.heartBeat = make(chan bool)   // 心跳信号的通道
	rf.heartbeatRe = make(chan bool) // 返回心跳的通道

	// 设计随机种子
	rand.Seed(time.Now().UnixNano())

	// 选举的协程
	go rf.election()
	// 心跳检测的协程
	go rf.sendLeaderHeartBeat()

	return rf

}

// 设置任期
func (rf *Raft) setTerm(term int) {
	rf.currentTerm = term
}

// 选举的方法
func (rf *Raft) election(){
	// 设置标记，判断是否选出了leader
	var result bool
	for {
		// 设置超时 150, 300 随机数
		timeout := randRange(150, 300)
		rf.lastMessageTime = millisecond()
		select {
		// 延迟等待一毫秒
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			fmt.Println("当前节点状态为：",rf.state)
		}
		result = false
		for !result {
			// 选主逻辑
			result = rf.electionOneRound(&leader)

		}
	}
}
// 随机值方法
func randRange(min,max int64) int64{
	return rand.Int63n(max-min) + min
}
// 获取当前时间，发送最后一条数据的时间
func millisecond() int64{
	return time.Now().UnixNano() / int64(time.Millisecond)
}
// 实现选主的逻辑
func (rf *Raft)electionOneRound(leader *Leader) bool{
	// 定义超时
	var timeout int64
	timeout = 100
	// 投票数量
	var vote int
	// 定义是否开始心跳新号的产生
	var triggerHeartbeat bool
	// 时间
	last := millisecond()
	// 用于返回值
	success := false

	// 将当前节点变成cadidate
	rf.mu.Lock()
	//修改状态
	rf.becomeCandidate()
	rf.mu.Unlock()
	fmt.Println("start electing leader")
	for {
		// 遍历所有节点拉选票
		for i := 0 ;i < raftCount; i++{
			if i != rf.me {
				// 拉选票
				go func() {
					if leader.LeaderId < 0{
						// 设置投票
						rf.electCh <- true
					}
				}()
			}
		}
		// 设置投票数量
		vote = 1
		// 遍历节点
		for i := 0 ;i < raftCount; i++{
			// 计算投票的数量
			select {
			case ok := <-rf.electCh:
				if ok{
					vote++ // 投票数量加1
					success = vote > raftCount / 2 // 若选票个数，大于节点个数 / 2，则成功
					if success && !triggerHeartbeat{
						// 变化成主节点，选主成功
						//触发心跳检测
						triggerHeartbeat = true
						rf.mu.Lock()
						rf.becomeLeader()	// 变主节点
						rf.mu.Unlock()
						// 由leader向其他节点发送心跳信号
						rf.heartBeat <- true
						fmt.Println(rf.me,"号节点称为了leader")
						fmt.Println("leader开始发送心跳信号了")
					}
				}
			}
		}
		// 做最后校验工作
		// 若不超时，且票数大于一半，则选举成功 break
		if timeout + last < millisecond() || (vote > raftCount / 2 || rf.currentLeader > -1 ){
			break
		}else {
			select {
				// 等待
				case <-time.After(time.Duration(10)*time.Millisecond):
			}
		}
	}
return success
}
// 修改状态candidate
func (rf *Raft)becomeCandidate(){
	rf.state = 1	// 将当前状态修改为1
	rf.setTerm(rf.currentTerm+1)	// 设置任期
	rf.votedFor = rf.me	// 给自己投票
	rf.currentLeader = -1
}
// 修改状态leader
func (rf *Raft)becomeLeader(){
	rf.state = 2
	rf.currentLeader = rf.me
}
// leader节点发送心跳信号
// 顺便完成数据同步 ，暂时未实现
// 看小弟挂没挂
func (rf *Raft)sendLeaderHeartBeat(){
	// 死循环
	for {
		select {
		case <-rf.heartBeat:
			rf.sendAppendEntriesImpl()
		}
	}
}
// 用于返回给leader的确认信号
func (rf *Raft)sendAppendEntriesImpl(){
	// 是主就别跑了
	if rf.currentLeader == rf.me{
		// 此时是leader
		var successCount = 0 // 记录确认信号的节点个数
		// 设置确认信号
		for i := 0 ; i < raftCount;i++{
			if i != rf.me{
				go func() {
					//rf.heartbeatRe <- true
					// 这里相当于客户端
					rp, err := rpc.DialHTTP("tcp", "127.0.0.1:8080")
					if err != nil{
						log.Fatal(err)
					}
					// 接收服务器返回的信息
					// 接收服务端返回信息的变量
					var ok = false
					err = rp.Call("Raft.Communication", Param{"hello"}, &ok)
					if err != nil {
						log.Fatal(err)
					}
					if ok {
						rf.heartbeatRe <- true
					}
				}()
			}
		}
		// 计算返回确认信号个数
		for i := 0;i < raftCount; i++{
			select {
			case ok:= <- rf.heartbeatRe:
				if ok {
					successCount++
					if successCount > raftCount / 2{
						fmt.Println("投票选举成功，心跳信号OK")
						log.Fatal("程序结束")
					}
				}
				
			}
		}
	}
}



// 首字母大写，RPC规范
// 分布式通信
type Param struct {
	Msg string
}

// 通信方法
func (r *Raft)Communication(p Param , a *bool) error{
	fmt.Println(p.Msg)
	*a  = true
	return nil
}