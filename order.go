package hftorderbook

// Single Order in an order book, as a node in a LimitOrder FIFO queue
type Order[P, V number] struct {
	Id       int
	Volume   V
	Next     *Order[P, V]
	Prev     *Order[P, V]
	Limit    *LimitOrder[P, V]
	BidOrAsk bool
}
