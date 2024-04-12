package user_store

import "cloud-time-tracker/cmd/officetime/api"

type userStore struct {
	user     api.User
	inOffice bool
}

type Store struct {
	users map[string]*userStore
}

func NewStore(users api.Users) *Store {
	m := map[string]*userStore{}
	for _, u := range users.Users {
		m[u.MacAddress] = &userStore{u, false}
	}
	return &Store{
		users: m,
	}
}

func (s *Store) setInOffice(macAddress string) {
	s.users[macAddress].inOffice = true
}
