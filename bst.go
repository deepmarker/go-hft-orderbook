package hftorderbook

import (
	"fmt"
)

// Simple Binary Search Tree, not self-balancing, good for random input

type nodeBST[P, V number] struct {
	Key   P
	Value *LimitOrder[P, V]
	Next  *nodeBST[P, V]
	Prev  *nodeBST[P, V]

	left  *nodeBST[P, V]
	right *nodeBST[P, V]
	size  int
}

type bst[P, V number] struct {
	root *nodeBST[P, V]
	minC *nodeBST[P, V] // cached min/max keys for O(1) access
	maxC *nodeBST[P, V]
}

func NewBST[P, V number]() bst[P, V] {
	return bst[P, V]{}
}

func (t *bst[_, _]) Size() int {
	return t.size(t.root)
}

func (t *bst[P, V]) size(n *nodeBST[P, V]) int {
	if n == nil {
		return 0
	}

	return n.size
}

func (t *bst[_, _]) IsEmpty() bool {
	return t.size(t.root) == 0
}

func (t *bst[_, _]) panicIfEmpty() {
	if t.IsEmpty() {
		panic("BST is empty")
	}
}

func (t *bst[P, _]) Contains(key P) bool {
	return t.get(t.root, key) != nil
}

func (t *bst[P, V]) Get(key P) *LimitOrder[P, V] {
	t.panicIfEmpty()

	x := t.get(t.root, key)
	if x == nil {
		panic(fmt.Sprintf("key %0.8f does not exist", key))
	}

	return x.Value
}

func (t *bst[P, V]) get(n *nodeBST[P, V], key P) *nodeBST[P, V] {
	if n == nil {
		return nil
	}

	if n.Key == key {
		return n
	}

	if n.Key > key {
		return t.get(n.left, key)
	} else {
		return t.get(n.right, key)
	}
}

func (t *bst[P, V]) Put(key P, value *LimitOrder[P, V]) {
	t.root = t.put(t.root, key, value)
}

func (t *bst[P, V]) put(n *nodeBST[P, V], key P, value *LimitOrder[P, V]) *nodeBST[P, V] {
	if n == nil {
		// search miss, creating a new node
		n := &nodeBST[P, V]{
			Value: value,
			Key:   key,
			size:  1,
		}

		if t.minC == nil || key < t.minC.Key {
			// new min
			t.minC = n
		}
		if t.maxC == nil || key > t.maxC.Key {
			// new max
			t.maxC = n
		}

		return n
	}

	if n.Key == key {
		// search hit, updating the value
		n.Value = value
		return n
	}

	if n.Key > key {
		left := n.left
		n.left = t.put(n.left, key, value)
		if left == nil {
			// new node has been just inserted to the left
			prev := n.Prev
			if prev != nil {
				prev.Next = n.left
			}
			n.left.Prev = prev
			n.left.Next = n
			n.Prev = n.left
		}
	} else {
		right := n.right
		n.right = t.put(n.right, key, value)
		if right == nil {
			// new node has been just inserted to the right
			next := n.Next
			if next != nil {
				next.Prev = n.right
			}
			n.right.Next = next
			n.right.Prev = n
			n.Next = n.right
		}
	}

	// re-calc size
	n.size = t.size(n.left) + 1 + t.size(n.right)
	return n
}

func (t *bst[P, V]) Height() int {
	if t.IsEmpty() {
		return 0
	}

	return t.height(t.root)
}

func (t *bst[P, V]) height(n *nodeBST[P, V]) int {
	if n == nil {
		return 0
	}

	lheight := t.height(n.left)
	rheight := t.height(n.right)

	height := lheight
	if rheight > lheight {
		height = rheight
	}

	return height + 1
}

func (t *bst[P, V]) Min() P {
	t.panicIfEmpty()
	return t.minC.Key
}

func (t *bst[P, V]) MinValue() *LimitOrder[P, V] {
	t.panicIfEmpty()
	return t.minC.Value
}

func (t *bst[P, V]) MinPointer() *nodeBST[P, V] {
	t.panicIfEmpty()
	return t.minC
}

func (t *bst[P, V]) min(n *nodeBST[P, V]) *nodeBST[P, V] {
	if n.left == nil {
		return n
	}

	return t.min(n.left)
}

func (t *bst[P, V]) Max() P {
	t.panicIfEmpty()
	return t.maxC.Key
}

func (t *bst[P, V]) MaxValue() *LimitOrder[P, V] {
	t.panicIfEmpty()
	return t.maxC.Value
}

func (t *bst[P, V]) MaxPointer() *nodeBST[P, V] {
	t.panicIfEmpty()
	return t.maxC
}

func (t *bst[P, V]) max(n *nodeBST[P, V]) *nodeBST[P, V] {
	if n.right == nil {
		return n
	}

	return t.max(n.right)
}

func (t *bst[P, V]) Floor(key P) P {
	t.panicIfEmpty()

	floor := t.floor(t.root, key)
	if floor == nil {
		panic(fmt.Sprintf("there are no keys <= %0.8f", key))
	}

	return floor.Key
}

func (t *bst[P, V]) floor(n *nodeBST[P, V], key P) *nodeBST[P, V] {
	if n == nil {
		// search miss
		return nil
	}

	if n.Key == key {
		// search hit
		return n
	}

	if n.Key > key {
		// floor must be in the left sub-tree
		return t.floor(n.left, key)
	}

	// key could be in the right sub-tree, if not, using current root
	floor := t.floor(n.right, key)
	if floor != nil {
		return floor
	}

	return n
}

func (t *bst[P, V]) Ceiling(key P) P {
	t.panicIfEmpty()

	ceiling := t.ceiling(t.root, key)
	if ceiling == nil {
		panic(fmt.Sprintf("there are no keys >= %0.8f", key))
	}

	return ceiling.Key
}

func (t *bst[P, V]) ceiling(n *nodeBST[P, V], key P) *nodeBST[P, V] {
	if n == nil {
		// search miss
		return nil
	}

	if n.Key == key {
		// search hit
		return n
	}

	if n.Key < key {
		// ceiling must be in the right sub-tree
		return t.ceiling(n.right, key)
	}

	// the key could be in the left sub-tree, if not, using current root
	ceiling := t.ceiling(n.left, key)
	if ceiling != nil {
		return ceiling
	}

	return n
}

func (t *bst[P, V]) Select(k int) P {
	if k < 0 || k >= t.Size() {
		panic("index out of range")
	}

	return t.selectNode(t.root, k).Key
}

func (t *bst[P, V]) selectNode(n *nodeBST[P, V], k int) *nodeBST[P, V] {
	if t.size(n.left) == k {
		return n
	}

	if t.size(n.left) > k {
		return t.selectNode(n.left, k)
	}

	k = k - t.size(n.left) - 1
	return t.selectNode(n.right, k)
}

func (t *bst[P, V]) Rank(key P) int {
	t.panicIfEmpty()
	return t.rank(t.root, key)
}

func (t *bst[P, V]) rank(n *nodeBST[P, V], key P) int {
	if n == nil {
		return 0
	}

	if n.Key == key {
		return t.size(n.left)
	}

	if n.Key > key {
		return t.rank(n.left, key)
	}

	return t.size(n.left) + 1 + t.rank(n.right, key)
}

func (t *bst[P, V]) deleteMin(n *nodeBST[P, V]) *nodeBST[P, V] {
	if n == nil {
		return nil
	}

	if n.left == nil {
		// we've reached the least leave of the tree
		next := n.Next
		prev := n.Prev
		if prev != nil {
			prev.Next = next
		}
		if next != nil {
			next.Prev = prev
		}
		n.Next = nil
		n.Prev = nil

		// updating global min
		if t.minC == n {
			t.minC = next
		}

		return n.right
	}

	n.left = t.deleteMin(n.left)

	// update size
	n.size = t.size(n.left) + 1 + t.size(n.right)
	return n
}

func (t *bst[P, V]) Delete(key P) {
	t.panicIfEmpty()

	t.root = t.delete(t.root, key)
}

func (t *bst[P, V]) delete(n *nodeBST[P, V], key P) *nodeBST[P, V] {
	if n == nil {
		return nil
	}

	if n.Key == key {
		// search hit

		// updating linked list
		next := n.Next
		prev := n.Prev
		if prev != nil {
			prev.Next = next
		}
		if next != nil {
			next.Prev = prev
		}
		n.Next = nil
		n.Prev = nil

		// updating global min and max
		if t.minC == n {
			t.minC = next
		}
		if t.maxC == n {
			t.maxC = prev
		}

		// replacing by successor (we can do similar with precedessor)
		if n.left == nil {
			return n.right
		} else if n.right == nil {
			return n.left
		}

		newn := t.min(n.right)
		newn.right = t.deleteMin(n.right)
		newn.left = n.left
		n = newn
	} else if n.Key > key {
		n.left = t.delete(n.left, key)
	} else {
		n.right = t.delete(n.right, key)
	}

	n.size = t.size(n.left) + 1 + t.size(n.right)
	return n
}

func (t *bst[P, V]) Keys(lo, hi P) []P {
	if lo < t.Min() || hi > t.Max() {
		panic("keys out of range")
	}

	return t.keys(t.root, lo, hi)
}

func (t *bst[P, V]) keys(n *nodeBST[P, V], lo, hi P) []P {
	if n == nil {
		return nil
	}

	if n.Key < lo {
		return t.keys(n.right, lo, hi)
	} else if n.Key > hi {
		return t.keys(n.left, lo, hi)
	}

	l := t.keys(n.left, lo, hi)
	r := t.keys(n.right, lo, hi)

	keys := make([]P, 0)
	if l != nil {
		keys = append(keys, l...)
	}
	keys = append(keys, n.Key)
	if r != nil {
		keys = append(keys, r...)
	}

	return keys
}

func (t *bst[P, V]) Print() {
	fmt.Println()
	t.print(t.root)
	fmt.Println()
}

func (t *bst[P, V]) print(n *nodeBST[P, V]) {
	if n == nil {
		return
	}

	fmt.Printf("%0.8f ", n.Key)

	t.print(n.left)
	t.print(n.right)
}
