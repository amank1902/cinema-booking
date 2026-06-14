package booking

import "sync"

type ConcurrentStore struct {
	bookings map[string]Booking
	sync.RWMutex
}

func NewConcurrentStore() *ConcurrentStore {
	return &ConcurrentStore{
		bookings: map[string]Booking{},
	}
}

func (m *ConcurrentStore) Book(b Booking) error {
	m.Lock()
	defer m.Unlock()
	if _, exists := m.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}
	m.bookings[b.SeatID] = b
	return nil
}

func (m *ConcurrentStore) ListBookings(movieId string) []Booking {
	m.RLock()
	defer m.RUnlock()
	
	var result []Booking
	for _, b := range m.bookings {
		if b.MovieID == movieId {
			result = append(result, b)
		}
	}
	return result
}