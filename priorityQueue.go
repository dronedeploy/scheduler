package main

import (
	"container/heap"
	"fmt"
	"strconv"
)

type PriorityKey struct {
	key      string
	priority int
	index    int
}

type KeyFunc func(obj interface{}) (string, error)

// MetaNamespaceKeyFunc is a convenient default KeyFunc which knows how to make
// keys for API objects
// TODO: add namespace to the name
func MetaNamespaceKeyFunc(obj interface{}) (string, error) {

	name := obj.(Pod).Metadata.Name
	return name, nil
}

type PriorityQueue struct {
	queue []PriorityKey
	items map[string]interface{}

	// keyFunc is used to make the key used for queued item insertion and retrieval, and
	// should be deterministic.
	keyFunc KeyFunc
}

// New PriorityQueue returns a struct that can be used to store items
// to be retrieved in priority order
func NewPriorityQueue(keyFunc KeyFunc) *PriorityQueue {
	pq := &PriorityQueue{
		items:   map[string]interface{}{},
		queue:   []PriorityKey{},
		keyFunc: keyFunc,
	}
	heap.Init(pq)
	return pq
}

// Methods for the heap interface:

// Len
// returns the length of the queue
func (pq PriorityQueue) Len() int {
	return len(pq.queue)
}

// Less
// compares the relative priority of two items in the queue and
// for a priority queue, less is more
func (pq PriorityQueue) Less(i, j int) bool {
	//Pop should give us the highest priority item
	return pq.more(i, j)
}

// more
func (pq PriorityQueue) more(i, j int) bool {
	return pq.queue[i].priority > pq.queue[j].priority
}

// Swap
// switches the position of two items in the queue
func (pq PriorityQueue) Swap(i, j int) {
	pq.queue[i], pq.queue[j] = pq.queue[j], pq.queue[i]
	pq.queue[i].index = i
	pq.queue[j].index = j
}

// Push
// adds an item to the queue. Used by the heap function. Do not use this
// outside heap.Push!
func (pq *PriorityQueue) Push(obj interface{}) {
	key, _ := pq.keyFunc(obj)
	priority, _ := MetaPriorityFunc(obj.(Pod))
	n := pq.Len()
	pk := PriorityKey{
		key:      key,
		priority: priority,
		index:    n,
	}

	pq.items[key] = obj
	pq.queue = append(pq.queue, pk)
}

// Pop
// grabs the highest priority item from the queue, returns and deletes it
// Do not use outside of heap.Pop!
func (pq *PriorityQueue) Pop() interface{} {
	//grab the queue
	old := pq.queue
	n := len(old)
	pk := old[n-1]
	item := pq.items[pk.key]

	//delete from map and array
	delete(pq.items, pk.key)
	pq.queue = old[0 : n-1]

	return item
}

// Helper Functions
// This is the annotation that determines the priority
const annotationKey = "k8s_priority"

// MetaPriorityFunc
// extracts the priority annotation of an object
// if the priority is not set, then set priority to -1
// The object must be a Pod
func MetaPriorityFunc(obj Pod) (int, error) {
	annotations := obj.Metadata.Annotations
	if p, ok := annotations[annotationKey]; ok {
		priority, err := strconv.Atoi(p)
		if err != nil {
			return -1, fmt.Errorf("priority is not an integer: %q", p)
		}
		return priority, nil
	}
	return -1, nil
}
