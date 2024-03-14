package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"6.5840/labgob"
	"bytes"
	"context"
	"fmt"
	//	"bytes"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	//	"6.5840/labgob"
	"6.5840/labrpc"
)

// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in part 2D you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh, but set CommandValid to false for these
// other uses.
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int

	// For 2D:
	SnapshotValid bool
	Snapshot      []byte
	SnapshotTerm  int
	SnapshotIndex int
}

const (
	FollowerState = iota
	CandidateState
	LeaderState
)

var (
	StateMap = map[int32]string{
		FollowerState:  "Follower",
		CandidateState: "Candidate",
		LeaderState:    "Leader",
	}
)

type Log struct {
	Term        int
	Index       int
	Command     interface{}
	ReceiveTime int64
	Committed   bool
}

type SyncEntryTask struct {
	Index int
	Done  chan struct{}
}

type SyncEntry struct {
	Once   sync.Once
	Server int
	Tasks  chan SyncEntryTask
}

// A Go object implementing a single Raft peer.
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	LeaderCtx    context.Context
	LeaderCancel context.CancelFunc

	applyCh       chan ApplyMsg
	Entries       []EntryArgs
	Logs          []Log
	NextIndex     []int
	SnapshotIndex int

	syncEntry []*SyncEntry

	VoteMaps map[int]int

	startTime    time.Time
	state        int32
	term         int
	voteTerm     int
	commitIndex  int
	commitMutex  sync.Mutex
	currentIndex int
}

func (rf *Raft) GetLogs(start, end int) []Log {
	return rf.Logs[start:end]
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	// Your code here (2A).
	return rf.getTerm(), rf.isState(LeaderState)
}

// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
// before you've implemented snapshots, you should pass nil as the
// second argument to persister.Save().
// after you've implemented snapshots, pass the current snapshot
// (or nil if there's not yet a snapshot).
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	if len(rf.Logs) == 0 {
		return
	}
	w := new(bytes.Buffer)
	e := labgob.NewEncoder(w)
	e.Encode(rf.term)
	e.Encode(rf.commitIndex)
	e.Encode(rf.currentIndex)
	e.Encode(rf.Logs)

	raftstate := w.Bytes()
	rf.persister.Save(raftstate, nil)

}

// restore previously persisted state.
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	r := bytes.NewBuffer(data)
	d := labgob.NewDecoder(r)

	var term int
	var commitIndex int
	var currentIndex int
	var logs []Log

	if d.Decode(&term) != nil || d.Decode(&commitIndex) != nil || d.Decode(&currentIndex) != nil || d.Decode(&logs) != nil {
		panic(fmt.Sprintf("readPersist error :%s", "decode error"))
	} else {
		rf.term = term
		rf.commitIndex = commitIndex
		rf.currentIndex = currentIndex
		rf.Logs = logs
	}
}

// the service says it has created a snapshot that has
// all info up to and including index. this means the
// service no longer needs the log through (and including)
// that index. Raft should now trim its log as much as possible.
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	// Your code here (2D).

}

// example RequestVote RPC arguments structure.
// field names must start with capital letters!
type RequestVoteArgs struct {
	SendTime int64
	Term     int
	From     int

	CurrentTerm  int
	CurrentIndex int

	ReceiveTime int64
	// Your data here (2A, 2B).
}

type EntryArgs struct {
	ReceiveTime int64
	Term        int

	LastLogIndex int
	LastLogTerm  int
	CommitIndex  int
	CommitTerm   int
	Logs         []Log
}

type EntryReply struct {
	IsSync bool
}

// example RequestVote RPC reply structure.
// field names must start with capital letters!
type RequestVoteReply struct {
	// Your data here (2A).

	VoteTerm int
	IsVote   bool
}

func (rf *Raft) AppendEntry(args *EntryArgs, reply *EntryReply) {
	if args.Term < rf.term {
		return
	}

	args.ReceiveTime = time.Now().UnixMilli()
	rf.mu.Lock()
	defer rf.mu.Unlock()

	rf.Entries = append(rf.Entries, *args)
	rf.setCommitted(args.CommitIndex, args.CommitTerm)
	if args.Term >= rf.term && !rf.isState(FollowerState) {
		rf.setState(FollowerState)
	}

	rf.term = args.Term
	reply.IsSync = false
	if len(args.Logs) > 0 {
		if args.LastLogIndex > 0 {
			if len(rf.Logs) <= args.LastLogIndex || rf.Logs[args.LastLogIndex].Term != args.LastLogTerm {
				return
			}

			rf.Logs = append(rf.Logs[:args.LastLogIndex+1], args.Logs...)
		} else {
			rf.Logs = args.Logs
			rf.SnapshotIndex = 0
		}

		reply.IsSync = true
	}
}

func (rf *Raft) GetLastLogIndexAndTerm() (int, int) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	return rf.GetLastLogIndexAndTermWithoutLock()
}

func (rf *Raft) GetLastLogIndexAndTermWithoutLock() (int, int) {
	if len(rf.Logs) == 0 {
		return 0, 0
	}

	return rf.Logs[len(rf.Logs)-1].Index, rf.Logs[len(rf.Logs)-1].Term
}

// example RequestVote RPC handler.
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if !rf.isState(FollowerState) && args.Term > rf.term {
		rf.term = args.Term
		rf.setState(FollowerState)
	}

	if args.Term <= rf.term || args.Term <= rf.voteTerm {
		reply.IsVote = false
		return
	}

	currentIndex, currentTerm := rf.GetLastLogIndexAndTermWithoutLock()
	if args.CurrentTerm < currentTerm || (args.CurrentTerm == currentTerm && args.CurrentIndex < currentIndex) {
		reply.IsVote = false
		return
	}

	if historyVote, ok := rf.VoteMaps[args.Term]; ok {
		reply.IsVote = false
		reply.VoteTerm = historyVote
		return
	}

	rf.VoteMaps[args.Term] = args.From
	rf.voteTerm = args.Term
	reply.IsVote = true
	reply.VoteTerm = args.Term
}

// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.

func (rf *Raft) sendAppendEntry(server int, args *EntryArgs, reply *EntryReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntry", args, reply)
	return ok
}

func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

func (rf *Raft) appendNewLog(command interface{}) int {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	rf.currentIndex += 1
	rf.Logs = append(rf.Logs, Log{
		Term:        rf.term,
		Command:     command,
		Index:       rf.currentIndex,
		ReceiveTime: time.Now().UnixMilli(),
	})

	return rf.currentIndex
}

func (rf *Raft) sendAppendEntryWithLog(server int, currentIndex int, done chan struct{}) {
	rf.syncEntry[server].Tasks <- SyncEntryTask{
		Index: currentIndex,
		Done:  done,
	}

	rf.syncEntry[server].Once.Do(func() {
		go func() {
			for {
				select {
				case <-rf.LeaderCtx.Done():
					return
				case task := <-rf.syncEntry[server].Tasks:
					args := &EntryArgs{
						Term:        rf.term,
						ReceiveTime: time.Now().UnixMilli(),
					}

					nextIndex := rf.getNextIndex(server)
					if nextIndex >= task.Index {
						continue
					}
					for {
						select {
						case <-rf.LeaderCtx.Done():
							return
						default:
						}

						if nextIndex >= 0 {
							args.LastLogTerm = rf.Logs[nextIndex].Term
							args.LastLogIndex = nextIndex
						} else {
							args.LastLogTerm = -1
							args.LastLogIndex = -1
							nextIndex = -1
						}

						args.Logs = rf.GetLogs(nextIndex, task.Index)

						reply := &EntryReply{}
						rf.sendAppendEntry(server, args, reply)
						if reply.IsSync {
							rf.NextIndex[server] = task.Index
							task.Done <- struct{}{}
							break
						}

						nextIndex -= 1
					}
				}
			}
		}()
	})
}

// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	term, isLeader := rf.GetState()
	if !isLeader {
		return -1, term, isLeader
	}

	currentIndex := rf.appendNewLog(command)
	go func() {
		done := make(chan struct{}, len(rf.peers))
		for i := 0; i < len(rf.peers); i++ {
			if i == rf.me {
				continue
			}

			rf.sendAppendEntryWithLog(i, currentIndex, done)
		}

		count := 1
		for rf.isState(LeaderState) {
			select {
			case <-done:
				count += 1
				if count > len(rf.peers)/2 && rf.isState(LeaderState) && term == rf.getTerm() {
					rf.setCommitted(currentIndex, rf.Logs[currentIndex].Term)
					return
				}
			case <-rf.LeaderCtx.Done():
				continue
			}
		}
	}()

	return currentIndex + 1, term, true
}

// the tester doesn't halt goroutines created by Raft after each test,
// but it does call the Kill() method. your code can use killed() to
// check whether Kill() has been called. the use of atomic avoids the
// need for a lock.
//
// the issue is that long-running goroutines use memory and may chew
// up CPU time, perhaps causing later tests to fail and generating
// confusing debug output. any goroutine with a long-running loop
// should call killed() to check whether it should stop.
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// Your code here, if desired.
	fmt.Printf("[%d] by killed\n", rf.me)
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}

func (rf *Raft) getTerm() int {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return rf.term
}

func (rf *Raft) setCommitted(index int, term int) {
	if index <= rf.commitIndex || len(rf.Logs) <= index || rf.Logs[index].Term != term {
		return
	}

	rf.commitMutex.Lock()
	defer rf.commitMutex.Unlock()

	if index <= rf.commitIndex || len(rf.Logs) <= index || rf.Logs[index].Term != term {
		return
	}

	for i := rf.commitIndex + 1; i <= index; i++ {
		rf.Logs[i].Committed = true

		fmt.Printf("[%d] apply index=%d, term=%d, command=%v\n", rf.me, rf.Logs[i].Index, rf.Logs[i].Term, rf.Logs[i].Command)

		rf.applyCh <- ApplyMsg{
			CommandValid: true,
			Command:      rf.Logs[i].Command,
			CommandIndex: rf.Logs[i].Index,
		}
	}

	rf.commitIndex = index
	rf.persist()
}

func (rf *Raft) setState(state int32) {
	fmt.Printf("[%d] %d set state %s, term=%d\n", rf.me, time.Now().Sub(rf.startTime).Milliseconds(), StateMap[state], rf.term)
	atomic.StoreInt32(&rf.state, state)
}

func (rf *Raft) isState(state int32) bool {
	return atomic.LoadInt32(&rf.state) == int32(state)
}

func (rf *Raft) getLastReceiveTime() (time.Time, bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if len(rf.Entries) == 0 {
		return time.Time{}, false
	}

	return time.UnixMilli(rf.Entries[len(rf.Entries)-1].ReceiveTime), true
}

func (rf *Raft) isFollowerTimeout(startTime time.Time) bool {
	if receiveTime, ok := rf.getLastReceiveTime(); ok {
		return receiveTime.Before(startTime)
	}

	return true
}

func (rf *Raft) startRequestVote() {
	rf.mu.Lock()
	rf.term++

	if _, ok := rf.VoteMaps[rf.term]; ok {
		rf.mu.Unlock()
		return
	}

	if rf.term <= rf.voteTerm {
		rf.mu.Unlock()
		return
	}

	rf.setState(CandidateState)
	rf.VoteMaps[rf.term] = rf.me

	rf.mu.Unlock()

	args := &RequestVoteArgs{
		Term:     rf.term,
		From:     rf.me,
		SendTime: time.Now().UnixMilli(),
	}

	if len(rf.Logs) > 0 {
		args.CurrentIndex, args.CurrentTerm = rf.GetLastLogIndexAndTerm()
	}

	replys := make([]*RequestVoteReply, len(rf.peers))

	var count int32 = 0
	done := make(chan struct{}, len(rf.peers))

	startTime := time.Now()
	for i := 0; i < len(rf.peers); i++ {
		if i == rf.me {
			continue
		}

		go func(server int) {
			defer func() {
				atomic.AddInt32(&count, 1)
				done <- struct{}{}
			}()

			reply := &RequestVoteReply{}
			rf.sendRequestVote(server, args, reply)

			replys[server] = reply
		}(i)
	}

	for {
		<-done

		voteTimes := 1
		for _, reply := range replys {
			if reply != nil && reply.IsVote {
				voteTimes += 1
			}
		}

		if voteTimes > (len(rf.peers)/2) || atomic.LoadInt32(&count)+1 == int32(len(rf.peers)) {
			fmt.Printf("[%d] start request vote success, term=%d, isCandidate=%v, voteTimes=%d, cost=%d\n",
				rf.me, rf.term, rf.isState(CandidateState), voteTimes, time.Now().Sub(startTime).Milliseconds())
			if voteTimes > (len(rf.peers)/2) && rf.isState(CandidateState) && args.Term == rf.getTerm() {
				rf.setState(LeaderState)
				go rf.SetLeader()
			}
			return
		}
	}
}

func (rf *Raft) findNextIndex() int {
	for i := len(rf.Logs) - 1; i >= 0; i-- {
		if rf.Logs[i].Committed {
			return rf.Logs[i].Index
		}
	}

	return 0
}

func (rf *Raft) getNextIndex(server int) int {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	return rf.NextIndex[server]
}

func (rf *Raft) SetLeader() {
	fmt.Printf("[%d] begin leader, term=%d, commitIndex=%d\n", rf.me, rf.term, rf.commitIndex)
	defer func() {
		rf.LeaderCancel()
		fmt.Printf("[%d] end leader\n", rf.me)
	}()

	rf.LeaderCtx, rf.LeaderCancel = context.WithCancel(context.Background())
	nextIndexList := make([]int, len(rf.peers))
	nextIndex := rf.findNextIndex()
	for i := 0; i < len(rf.peers); i++ {
		nextIndexList[i] = nextIndex
	}

	rf.NextIndex = nextIndexList

	rf.syncEntry = make([]*SyncEntry, len(rf.peers))
	for i := 0; i < len(rf.peers); i++ {
		rf.syncEntry[i] = &SyncEntry{
			Server: i,
			Once:   sync.Once{},
			Tasks:  make(chan SyncEntryTask, 100),
		}
	}

	rf.currentIndex = 0
	if len(rf.Logs) > 0 {
		rf.currentIndex = rf.Logs[len(rf.Logs)-1].Index
	}

	ticker := time.NewTicker(50 * time.Millisecond)
	for {
		<-ticker.C
		if rf.killed() || !rf.isState(LeaderState) {
			return
		}

		args := &EntryArgs{
			Term: rf.getTerm(),
		}

		if rf.commitIndex > 0 {
			args.CommitIndex = rf.commitIndex
			args.CommitTerm = rf.Logs[rf.commitIndex].Term
		}
		for i := 0; i < len(rf.peers); i++ {
			if i == rf.me {
				continue
			}

			go func(server int) {
				reply := &EntryReply{}
				rf.sendAppendEntry(server, args, reply)
			}(i)
		}
	}
}

func (rf *Raft) ticker() {
	for rf.killed() == false {
		startTime := time.Now()

		// Your code here (2A)
		// Check if a leader election should be started.

		// pause for a random amount of time between 50 and 350
		// milliseconds.
		ms := 50 + (rand.Int63() % 300)

		time.Sleep(time.Duration(ms) * time.Millisecond)

		if rf.isState(LeaderState) {
			continue
		}

		if rf.isState(FollowerState) {
			if rf.isFollowerTimeout(startTime) {
				go rf.startRequestVote()
			}
		} else if rf.isState(CandidateState) {
			go rf.startRequestVote()
		}
	}
}

// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.startTime = time.Now()
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	rf.setState(FollowerState)
	rf.VoteMaps = make(map[int]int)
	rf.commitIndex = 0
	rf.applyCh = applyCh
	rf.currentIndex = 0
	// Your initialization code here (2A, 2B, 2C).

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	// start ticker goroutine to start elections
	go rf.ticker()

	return rf
}

// TODO
// 将 raft 日志的 index 修改为 1 开始
