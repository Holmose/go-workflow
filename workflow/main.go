package main

import (
	"context"
	"fmt"
)

// WearUnderpantsAction 穿内裤任务
type WearUnderpantsAction struct{}

func (a *WearUnderpantsAction) Run(i interface{}) {
	fmt.Println("我正在穿内裤")
}

// WearSocksAction 穿袜子任务
type WearSocksAction struct{}

func (a *WearSocksAction) Run(i interface{}) {
	fmt.Println("我正在穿袜子")
}

// ShirtNodeAction 穿衬衣任务
type ShirtNodeAction struct{}

func (a *ShirtNodeAction) Run(i interface{}) {
	fmt.Println("我正在穿衬衣")
}

// WatchNodeAction 戴手表任务
type WatchNodeAction struct{}

func (a *WatchNodeAction) Run(i interface{}) {
	fmt.Println("我正在戴手表")
}

// WearTrouserNodeAction 穿裤子任务
type WearTrouserNodeAction struct{}

func (a *WearTrouserNodeAction) Run(i interface{}) {
	fmt.Println("我正在穿裤子")
}

// WearCoatNodeAction 穿外套任务
type WearCoatNodeAction struct{}

func (a *WearCoatNodeAction) Run(i interface{}) {
	fmt.Println("我正在穿外套")
}

// WearShoesNodeAction 穿鞋子任务
type WearShoesNodeAction struct{}

func (a *WearShoesNodeAction) Run(i interface{}) {
	fmt.Println("我正在穿鞋子")
}

func main() {
	wf := NewWorkFlow()
	// 构建节点
	UnderpantsNode := NewNode(&WearUnderpantsAction{}) // 内裤
	SocksNode := NewNode(&WearSocksAction{})           // 袜子
	ShirtNode := NewNode(&ShirtNodeAction{})           // 衬衣
	WatchNode := NewNode(&WatchNodeAction{})           // 手表

	TrousersNode := NewNode(&WearTrouserNodeAction{}) // 裤子
	ShoesNode := NewNode(&WearShoesNodeAction{})      // 鞋子
	CoatNode := NewNode(&WearCoatNodeAction{})        // 外套

	// 构建节点之间的关系
	// 启始节点
	wf.AddStartNode(UnderpantsNode)
	wf.AddStartNode(SocksNode)
	wf.AddStartNode(ShirtNode)
	wf.AddStartNode(WatchNode)

	// 中间节点
	wf.AddEdge(UnderpantsNode, TrousersNode)
	wf.AddEdge(TrousersNode, ShoesNode)
	wf.AddEdge(SocksNode, ShoesNode)

	wf.AddEdge(ShirtNode, CoatNode)
	wf.AddEdge(WatchNode, CoatNode)

	// 收尾节点
	wf.ConnectToEnd(ShoesNode)
	wf.ConnectToEnd(CoatNode)

	// 数据
	var completedAction []string

	ctx, _ := context.WithCancel(context.Background())
	wf.StartWithContext(ctx, completedAction)
	wf.WaitDone()

	fmt.Println("执行其他逻辑")
}
