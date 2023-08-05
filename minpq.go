package hftorderbook

// Mininum oriented Priority Queue
type minPQ[P number] struct {
	keys []P
	n    int
}

func NewMinPQ[P number](size int) minPQ[P] {
	return minPQ[P]{
		keys: make([]P, size+1),
	}
}

func (pq *minPQ[P]) Size() int {
	return pq.n
}

func (pq *minPQ[P]) IsEmpty() bool {
	return pq.n == 0
}

func (pq *minPQ[P]) Insert(key P) {
	if pq.n+1 == cap(pq.keys) {
		panic("pq is full")
	}

	pq.n++
	pq.keys[pq.n] = key

	// restore order: LogN
	pq.swim(pq.n)
}

func (pq *minPQ[P]) Top() P {
	if pq.IsEmpty() {
		panic("pq is empty")
	}

	return pq.keys[1]
}

// removes minimal element and returns it
func (pq *minPQ[P]) DelTop() P {
	if pq.IsEmpty() {
		panic("pq is empty")
	}

	top := pq.keys[1]
	pq.keys[1] = pq.keys[pq.n]
	pq.n--

	// restore order: logN
	pq.sink(1)

	return top
}

func (pq *minPQ[P]) swim(k int) {
	for k > 1 && pq.keys[k] < pq.keys[k/2] {
		// swap
		pq.keys[k], pq.keys[k/2] = pq.keys[k/2], pq.keys[k]
		k = k / 2
	}
}

func (pq *minPQ[P]) sink(k int) {
	for 2*k <= pq.n {
		c := 2 * k
		// select minimum of two children
		if c < pq.n && pq.keys[c+1] < pq.keys[c] {
			c++
		}

		if pq.keys[c] < pq.keys[k] {
			// swap
			pq.keys[c], pq.keys[k] = pq.keys[k], pq.keys[c]
			k = c
		} else {
			break
		}
	}
}
