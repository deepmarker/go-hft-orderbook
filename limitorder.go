package hftorderbook

// Limit price orders combined as a FIFO queue
type LimitOrder[P, V number] struct {
	Price P

	orders      *ordersQueue[P, V]
	totalVolume V
}

func NewLimitOrder[P, V number](price P) LimitOrder[P, V] {
	q := ordersQueue[P, V]{}
	return LimitOrder[P, V]{
		Price:  price,
		orders: &q,
	}
}

func (this *LimitOrder[_, V]) TotalVolume() V {
	return this.totalVolume
}

func (this *LimitOrder[_, _]) Size() int {
	return this.orders.Size()
}

func (this *LimitOrder[P, V]) Enqueue(o *Order[P, V]) {
	this.orders.Enqueue(o)
	o.Limit = this
	this.totalVolume += o.Volume
}

func (this *LimitOrder[P, V]) Dequeue() *Order[P, V] {
	if this.orders.IsEmpty() {
		return nil
	}

	o := this.orders.Dequeue()
	this.totalVolume -= o.Volume
	return o
}

func (this *LimitOrder[P, V]) Delete(o *Order[P, V]) {
	if o.Limit != this {
		panic("order does not belong to the limit")
	}

	this.orders.Delete(o)
	o.Limit = nil
	this.totalVolume -= o.Volume
}

func (this *LimitOrder[P, V]) Clear() {
	q := NewOrdersQueue[P, V]()
	this.orders = &q
	this.totalVolume = 0
}
