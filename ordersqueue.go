package hftorderbook

// Doubly linked orders queue
// TODO: this should be compared with ring buffer queue performance
type ordersQueue[P, V number] struct {
	head *Order[P, V]
	tail *Order[P, V]
	size int
}

func NewOrdersQueue[P, V number]() ordersQueue[P, V] {
	return ordersQueue[P, V]{}
}

func (this *ordersQueue[_, _]) Size() int {
	return this.size
}

func (this *ordersQueue[_, _]) IsEmpty() bool {
	return this.size == 0
}

func (this *ordersQueue[P, V]) Enqueue(o *Order[P, V]) {
	tail := this.tail
	this.tail = o
	if tail != nil {
		tail.Next = o
		o.Prev = tail
	}
	if this.head == nil {
		this.head = o
	}
	this.size++
}

func (this *ordersQueue[P, V]) Dequeue() *Order[P, V] {
	if this.size == 0 {
		return nil
	}

	head := this.head
	if this.tail == this.head {
		this.tail = nil
	}

	this.head = this.head.Next
	this.size--
	return head
}

func (this *ordersQueue[P, V]) Delete(o *Order[P, V]) {
	prev := o.Prev
	next := o.Next
	if prev != nil {
		prev.Next = next
	}
	if next != nil {
		next.Prev = prev
	}
	o.Next = nil
	o.Prev = nil

	this.size--

	if this.head == o {
		this.head = next
	}
	if this.tail == o {
		this.tail = prev
	}
}
