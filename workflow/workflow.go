package workflow

// Go基于有向无环图的并发执行流
// 参考链接：https://cc.topgoer.cn/blog-230.html

import (
	"context"
	"sync"
	"sync/atomic"
)

type WorkFlow struct {
	done        chan struct{} // 结束标识,该标识由结束节点写入
	doneOnce    *sync.Once    // 保证并发时只写入一次
	alreadyDone bool          // 有节点出错时终止流程标记
	root        *Node         // 开始节点
	End         *Node         // 结束节点
	edges       []*Edge       // 所有经过的边，边连接了节点
}
type Edge struct {
	FromNode *Node // 前置节点
	ToNode   *Node // 下一节点
}
type Node struct {
	Dependency   []*Edge  // 依赖的边
	DepCompleted int32    // 表示依赖的边有多少个已执行完成，用于判断该节点是否可以执行了
	Task         Runnable // 任务执行
	Children     []*Edge  // 节点的孩子边
}

// Runnable 对节点执行业务逻辑的抽象 i:为数据，可在函数中进行调用
type Runnable interface {
	Run(i interface{})
}

func NewNode(Task Runnable) *Node {
	return &Node{
		Task: Task,
	}
}

// AddEdge 将两个节点构建成边
func AddEdge(from *Node, to *Node) *Edge {
	edg := &Edge{
		FromNode: from,
		ToNode:   to,
	}
	// 该条边是from节点的出边
	from.Children = append(from.Children, edg)
	// 该条边是to节点的入边
	to.Dependency = append(to.Dependency, edg)
	return edg
}

// EndWorkFlowAction 结束节点的具体执行任务
type EndWorkFlowAction struct {
	done chan struct{} // 节点执行完成，往该done写入消息，和workflow中的done共用
	s    *sync.Once    // 并发控制，确保只往done中写入一次
}

func (end *EndWorkFlowAction) Run(i interface{}) {
	end.s.Do(func() {
		end.done <- struct{}{}
	})
}

func NewWorkFlow() *WorkFlow {
	wf := &WorkFlow{
		root:     &Node{Task: nil}, // 开始节点，所有具体的节点都是它的子节点，没有具体的执行逻辑，只为出发其他节点的执行
		done:     make(chan struct{}, 1),
		doneOnce: &sync.Once{},
	}

	// 加入结束节点
	EndNode := &EndWorkFlowAction{
		done: wf.done,
		s:    wf.doneOnce,
	}
	wf.End = NewNode(EndNode)
	return wf
}

func (wf *WorkFlow) AddStartNode(node *Node) {
	// 根节点和中间节点之间构成的边
	wf.edges = append(wf.edges, AddEdge(wf.root, node))
}
func (wf *WorkFlow) AddEdge(from *Node, to *Node) {
	// 中间节点和中间节点之间构成的边
	wf.edges = append(wf.edges, AddEdge(from, to))
}
func (wf *WorkFlow) ConnectToEnd(node *Node) {
	// 中间节点和结束节点之间构成的边
	wf.edges = append(wf.edges, AddEdge(node, wf.End))
}

func (wf *WorkFlow) StartWithContext(ctx context.Context, i interface{}) {
	wf.root.ExecuteWithContext(ctx, wf, i)
}
func (wf *WorkFlow) WaitDone() {
	<-wf.done
	close(wf.done)
}
func (wf *WorkFlow) interruptDone() {
	wf.alreadyDone = true
	wf.doneOnce.Do(func() {
		wf.done <- struct{}{}
	})
}

// ExecuteWithContext 节点的具体执行逻辑实现
func (n *Node) ExecuteWithContext(ctx context.Context, wf *WorkFlow, i interface{}) {
	// 所依赖的前置节点没有运行完成，则直接返回
	if !n.dependencyHasDone() {
		return
	}
	// 有节点运行出错，终止流程的执行
	if ctx.Err() != nil {
		wf.interruptDone()
		return
	}
	// 节点具体的运行逻辑
	if n.Task != nil {
		n.Task.Run(i)
	}

	// 运行子节点
	if len(n.Children) > 0 {
		for idx := 1; idx < len(n.Children); idx++ {
			go func(child *Edge) {
				child.ToNode.ExecuteWithContext(ctx, wf, i)
			}(n.Children[idx])
		}

		n.Children[0].ToNode.ExecuteWithContext(ctx, wf, i)
	}
}

func (n *Node) dependencyHasDone() bool {
	// 该节点没有依赖的前置节点，不需要等待，直接返回true
	if n.Dependency == nil {
		return true
	}

	// 如果该节点只有一个依赖的前置节点，也直接返回
	if len(n.Dependency) == 1 {
		return true
	}

	// 这里将依赖的节点加1，说明有一个依赖的节点完成了
	atomic.AddInt32(&n.DepCompleted, 1)

	// 判断当前依赖的节点数量是否和依赖的节点相等，相等，说明都运行完了
	return n.DepCompleted == int32(len(n.Dependency))
}
