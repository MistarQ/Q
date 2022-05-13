package core

import (
	"container/heap"
	"fmt"
	"sort"
	"testing"
)

//type maxHeap []int // 大顶堆
//type minHeap []int // 小顶堆
//
//// 每个堆都要heap.Interface的五个方法：Len, Less, Swap, Push, Pop
//// 其实只有Less的区别。
//
//// Len 返回堆的大小
//func (m maxHeap) Len() int {
//	return len(m)
//}
//func (m minHeap) Len() int {
//	return len(m)
//}
//
//// Less 决定是大优先还是小优先
//func (m maxHeap) Less(i, j int) bool { // 大根堆
//	return m[i] > m[j]
//}
//func (m minHeap) Less(i, j int) bool { // 小根堆
//	return m[i] < m[j]
//}
//
//// Swap 交换下标i, j元素的顺序
//func (m maxHeap) Swap(i, j int) {
//	m[i], m[j] = m[j], m[i]
//}
//func (m minHeap) Swap(i, j int) {
//	m[i], m[j] = m[j], m[i]
//}
//
//// Push 在堆的末尾添加一个元素，注意和heap.Push(heap.Interface, interface{})区分
//func (m *maxHeap) Push(v interface{}) {
//	*m = append(*m, v.(int))
//}
//func (m *minHeap) Push(v interface{}) {
//	*m = append(*m, v.(int))
//}
//
//// Pop 删除堆尾的元素，注意和heap.Pop()区分
//func (m *maxHeap) Pop() interface{} {
//	old := *m
//	n := len(old)
//	v := old[n-1]
//	*m = old[:n-1]
//	return v
//}
//func (m *minHeap) Pop() interface{} {
//	old := *m
//	n := len(old)
//	v := old[n-1]
//	*m = old[:n-1]
//	return v
//}
//
//// MedianFinder 维护两个堆，前一半是大顶堆，后一半是小顶堆，中位数由两个堆顶决定
//type MedianFinder struct {
//	maxH *maxHeap
//	minH *minHeap
//}
//
//// Constructor 建两个空堆
//func Constructor() MedianFinder {
//	return MedianFinder{
//		new(maxHeap),
//		new(minHeap),
//	}
//}
//
//// AddNum 插入元素num
//// 分两种情况插入：
//// 1. 两个堆的大小相等，则小顶堆增加一个元素（增加的不一定是num）
//// 2. 小顶堆比大顶堆多一个元素，大顶堆增加一个元素
//// 这两种情况又分别对应两种情况：
//// 1. num小于大顶堆的堆顶，则num插入大顶堆
//// 2. num大于小顶堆的堆顶，则num插入小顶堆
//// 插入完成后记得调整堆的大小使得两个堆的容量相等，或小顶堆大1
//func (m *MedianFinder) AddNum(num int) {
//	if m.maxH.Len() == m.minH.Len() {
//		if m.minH.Len() == 0 || num >= (*m.minH)[0] {
//			heap.Push(m.minH, num)
//		} else {
//			heap.Push(m.maxH, num)
//			top := heap.Pop(m.maxH).(int)
//			heap.Push(m.minH, top)
//		}
//	} else {
//		if num > (*m.minH)[0] {
//			heap.Push(m.minH, num)
//			bottle := heap.Pop(m.minH).(int)
//			heap.Push(m.maxH, bottle)
//		} else {
//			heap.Push(m.maxH, num)
//		}
//	}
//}
//
//// FindMediam 输出中位数
//func (m *MedianFinder) FindMedian() float64 {
//	if m.minH.Len() == m.maxH.Len() {
//		return float64((*m.maxH)[0])/2.0 + float64((*m.minH)[0])/2.0
//	} else {
//		return float64((*m.minH)[0])
//	}
//}

/**
 * Your MedianFinder object will be instantiated and called as such:
 * obj := Constructor();
 * obj.AddNum(num);
 * param_2 := obj.FindMedian();
 */

type MedianFinder struct {
	queMin, queMax hp
}

func Constructor() MedianFinder {
	return MedianFinder{}
}

func (mf *MedianFinder) AddNum(num int) {
	minQ, maxQ := &mf.queMin, &mf.queMax
	if minQ.Len() == 0 || num <= -minQ.IntSlice[0] {
		heap.Push(minQ, -num)
		if maxQ.Len()+1 < minQ.Len() {
			heap.Push(maxQ, -heap.Pop(minQ).(int))
		}
	} else {
		heap.Push(maxQ, num)
		if maxQ.Len() > minQ.Len() {
			heap.Push(minQ, -heap.Pop(maxQ).(int))
		}
	}
}

func (mf *MedianFinder) FindMedian() float64 {
	minQ, maxQ := mf.queMin, mf.queMax
	if minQ.Len() > maxQ.Len() {
		return float64(-minQ.IntSlice[0])
	}
	return float64(maxQ.IntSlice[0]-minQ.IntSlice[0]) / 2
}

type hp struct{ sort.IntSlice }

func (h *hp) Push(v interface{}) { h.IntSlice = append(h.IntSlice, v.(int)) }
func (h *hp) Pop() (v interface{}) {
	a := h.IntSlice
	v = a[len(a)-1]
	h.IntSlice = a[:len(a)-1]
	return v
}

func TestSpiralOrder(t *testing.T) {
	arr := [][]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
	}
	fmt.Println(spiralOrder(arr))
}

var visited [][]bool
var ans []int
var boundX int
var boundY int

func spiralOrder(matrix [][]int) []int {
	boundX = len(matrix[0]) - 1
	boundY = len(matrix) - 1
	for _, v := range matrix {
		arr := make([]bool, len(v))
		visited = append(visited, arr)
	}

	x := 0
	y := 0
	for {
		if !visited[y][x] {
			ans = append(ans, matrix[y][x])
			visited[y][x] = true
		}
		if x < boundX && !visited[y][x+1] {
			x++
		} else if y < boundY && !visited[y+1][x] {
			y++

		} else if x > 0 && !visited[y][x-1] {
			x--
		} else if y > 0 && !visited[y-1][x] {
			y--
		} else {
			break
		}
	}
	return ans
}
