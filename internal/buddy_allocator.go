package internal

import "math/bits"

type BuddyAllocator struct {
	capacity int
	maxOrder int

	freeSet   []map[int]struct{}
	freeStack [][]int
}

func NewBuddyAllocator(initialCapacity int, strict bool) *BuddyAllocator {
	if initialCapacity < 1 {
		initialCapacity = 1
	}

	capPow2 := nextPow2(initialCapacity)
	maxOrder := ilog2(capPow2)

	buddyAllocator := &BuddyAllocator{
		capacity:  capPow2,
		maxOrder:  maxOrder,
		freeSet:   make([]map[int]struct{}, maxOrder+1),
		freeStack: make([][]int, maxOrder+1),
	}

	for i := range buddyAllocator.freeSet {
		buddyAllocator.freeSet[i] = make(map[int]struct{}, 16)
	}

	buddyAllocator.pushFree(maxOrder, 0)
	return buddyAllocator
}

func (this *BuddyAllocator) Allocate(size int) (index int, ok bool) {
	if size <= 0 {
		panic("Allocate: size must be > 0")
	}

	blockSize := nextPow2(size)
	if blockSize > this.capacity {
		return 0, false
	}

	reqOrder := ilog2(blockSize)

	for order := reqOrder; order <= this.maxOrder; order++ {
		var got bool
		index, got = this.popFree(order)
		if !got {
			continue
		}

		for order > reqOrder {
			order--
			half := 1 << order
			rightBuddy := index + half
			this.pushFree(order, rightBuddy)
		}

		return index, true
	}

	return 0, false
}

func (this *BuddyAllocator) Capacity() int {
	return this.capacity
}

func (this *BuddyAllocator) Free(index int, size int) {
	if size <= 0 {
		panic("Free: size must be > 0")
	}

	blockSize := nextPow2(size)
	if blockSize > this.capacity {
		panic("Free: size exceeds allocator capacity")
	}

	order := ilog2(blockSize)

	if index < 0 || index+blockSize > this.capacity {
		panic("Free: index out of range")
	}

	if index&(blockSize-1) != 0 {
		panic("Free: index not aligned to block size")
	}

	for order < this.maxOrder {
		buddy := index ^ blockSize // buddy formula for buddy allocator
		if this.isFree(order, buddy) {
			this.removeFree(order, buddy)
			if buddy < index {
				index = buddy
			}

			order++
			blockSize <<= 1
			continue
		}

		break
	}

	this.pushFree(order, index)
}

func (this *BuddyAllocator) Grow() (newSize int) {
	target := this.capacity * 2

	oldCap := this.capacity
	oldMax := this.maxOrder

	newMax := ilog2(target)
	if newMax > oldMax {
		this.freeSet = append(this.freeSet, make([]map[int]struct{}, newMax-oldMax)...)
		this.freeStack = append(this.freeStack, make([][]int, newMax-oldMax)...)
		for i := oldMax + 1; i <= newMax; i++ {
			this.freeSet[i] = make(map[int]struct{}, 16)
		}
	}

	this.capacity = target
	this.maxOrder = newMax

	for stepCap := oldCap; stepCap < target; stepCap <<= 1 {
		this.Free(stepCap, stepCap)
	}

	return this.capacity
}

// --- private methods ---

func (this *BuddyAllocator) isFree(order int, index int) bool {
	_, ok := this.freeSet[order][index]
	return ok
}

func (this *BuddyAllocator) popFree(order int) (index int, ok bool) {
	stack := this.freeStack[order]
	for len(stack) > 0 {
		n := len(stack) - 1
		index = stack[n]
		stack = stack[:n]

		if _, exists := this.freeSet[order][index]; !exists {
			continue
		}

		delete(this.freeSet[order], index)
		this.freeStack[order] = stack
		return index, true
	}

	this.freeStack[order] = stack
	return 0, false
}

func (this *BuddyAllocator) pushFree(order int, index int) {
	if _, exists := this.freeSet[order][index]; exists {
		panic("pushFree: double-free detected")
	}

	this.freeSet[order][index] = struct{}{}
	this.freeStack[order] = append(this.freeStack[order], index)
}

func (this *BuddyAllocator) removeFree(order int, index int) {
	if _, ok := this.freeSet[order][index]; !ok {
		panic("removeFree: not free")
	}

	delete(this.freeSet[order], index)
}

// --- private functions ---

func ilog2(pow2 int) int {
	if pow2 <= 0 || (pow2&(pow2-1)) != 0 {
		panic("ilog2: input must be power of two")
	}

	return bits.Len(uint(pow2)) - 1
}

func nextPow2(n int) int {
	if n <= 1 {
		return 1
	}

	return 1 << bits.Len(uint(n-1))
}
