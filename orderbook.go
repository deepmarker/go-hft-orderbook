package hftorderbook

import (
	"fmt"
	"sync"
)

// maximum limits per orderbook side to pre-allocate memory
const MaxLimitsNum int = 10000

type Orderbook[P, V number] struct {
	Bids *redBlackBST[P, V]
	Asks *redBlackBST[P, V]

	bidLimitsCache map[P]*LimitOrder[P, V]
	askLimitsCache map[P]*LimitOrder[P, V]
	pool           *sync.Pool
}

func NewOrderbook[P, V number]() Orderbook[P, V] {
	bids := NewRedBlackBST[P, V]()
	asks := NewRedBlackBST[P, V]()
	return Orderbook[P, V]{
		Bids: &bids,
		Asks: &asks,

		bidLimitsCache: make(map[P]*LimitOrder[P, V], MaxLimitsNum),
		askLimitsCache: make(map[P]*LimitOrder[P, V], MaxLimitsNum),
		pool: &sync.Pool{
			New: func() interface{} {
				limit := NewLimitOrder[P, V](0)
				return &limit
			},
		},
	}
}

func (this *Orderbook[P, V]) Add(price P, o *Order[P, V]) {
	var limit *LimitOrder[P, V]

	if o.BidOrAsk {
		limit = this.bidLimitsCache[price]
	} else {
		limit = this.askLimitsCache[price]
	}

	if limit == nil {
		// getting a new limit from pool
		limit = this.pool.Get().(*LimitOrder[P, V])
		limit.Price = price

		// insert into the corresponding BST and cache
		if o.BidOrAsk {
			this.Bids.Put(price, limit)
			this.bidLimitsCache[price] = limit
		} else {
			this.Asks.Put(price, limit)
			this.askLimitsCache[price] = limit
		}
	}

	// add order to the limit
	limit.Enqueue(o)
}

func (this *Orderbook[P, V]) Cancel(o *Order[P, V]) {
	limit := o.Limit
	limit.Delete(o)

	if limit.Size() == 0 {
		// remove the limit if there are no orders
		if o.BidOrAsk {
			this.Bids.Delete(limit.Price)
			delete(this.bidLimitsCache, limit.Price)
		} else {
			this.Asks.Delete(limit.Price)
			delete(this.askLimitsCache, limit.Price)
		}

		// put it back to the pool
		this.pool.Put(limit)
	}
}

func (this *Orderbook[P, V]) ClearBidLimit(price P) {
	this.clearLimit(price, true)
}

func (this *Orderbook[P, V]) ClearAskLimit(price P) {
	this.clearLimit(price, false)
}

func (this *Orderbook[P, V]) clearLimit(price P, bidOrAsk bool) {
	var limit *LimitOrder[P, V]
	if bidOrAsk {
		limit = this.bidLimitsCache[price]
	} else {
		limit = this.askLimitsCache[price]
	}

	if limit == nil {
		panic(fmt.Sprintf("there is no such price limit %0.8f", price))
	}

	limit.Clear()
}

func (this *Orderbook[P, V]) DeleteBidLimit(price P) {
	limit := this.bidLimitsCache[price]
	if limit == nil {
		return
	}

	this.deleteLimit(price, true)
	delete(this.bidLimitsCache, price)

	// put limit back to the pool
	limit.Clear()
	this.pool.Put(limit)

}

func (this *Orderbook[P, V]) DeleteAskLimit(price P) {
	limit := this.askLimitsCache[price]
	if limit == nil {
		return
	}

	this.deleteLimit(price, false)
	delete(this.askLimitsCache, price)

	// put limit back to the pool
	limit.Clear()
	this.pool.Put(limit)
}

func (this *Orderbook[P, V]) deleteLimit(price P, bidOrAsk bool) {
	if bidOrAsk {
		this.Bids.Delete(price)
	} else {
		this.Asks.Delete(price)
	}
}

func (this *Orderbook[P, V]) GetVolumeAtBidLimit(price P) *V {
	limit := this.bidLimitsCache[price]
	if limit == nil {
		return nil
	}
	v := limit.TotalVolume()
	return &v
}

func (this *Orderbook[P, V]) GetVolumeAtAskLimit(price P) *V {
	limit := this.askLimitsCache[price]
	if limit == nil {
		return nil
	}
	v := limit.TotalVolume()
	return &v
}

func (this *Orderbook[P, V]) GetBestBid() P {
	return this.Bids.Max()
}

func (this *Orderbook[P, V]) GetBestOffer() P {
	return this.Asks.Min()
}

func (this *Orderbook[P, V]) BLength() int {
	return len(this.bidLimitsCache)
}

func (this *Orderbook[P, V]) ALength() int {
	return len(this.askLimitsCache)
}
