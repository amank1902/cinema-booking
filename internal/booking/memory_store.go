package booking

type MemoryStore struct{
	bookings map[string]Booking
}

func NewMemoryStore() *MemoryStore{
	return &MemoryStore{
		bookings: map[string]Booking{},
	}
}

func (m *MemoryStore) Book(b Booking) error{
	if _,exists := m.bookings[b.SeatID]; exists{
		return ErrSeatAlreadyBooked
	}
	m.bookings[b.SeatID] = b
	return nil
}

func (m *MemoryStore)ListBookings(movieId string) []Booking{
	var result []Booking
	for _,b := range m.bookings{
		if b.MovieID == movieId{
			result = append(result, b)
		}
	}
	return result
}