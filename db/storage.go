package chat_db

import (
	"sort"
	chat_type "web_server/type"
)

type Handler interface {
	Set(key string, value any, group ...string) error
	Get(key string, value any, group ...string) (bool, error)
}

type RoomStorage struct {
	handler Handler
	rooms   chat_type.RoomList
}

func (s *RoomStorage) LoadAll() (chat_type.RoomList, error) {
	if len(s.rooms) > 0 {
		return s.rooms, nil
	}
	records := chat_type.RoomList{}
	_, err := s.handler.Get("room_list", &records)
	sort.Sort(sort.StringSlice(records))
	s.rooms = records
	return records, err
}

func (s *RoomStorage) SaveAll(rooms chat_type.RoomList) error {
	sort.Sort(sort.StringSlice(rooms))
	s.rooms = rooms
	return s.handler.Set("room_list", rooms)
}

func (s *RoomStorage) FindRoom(room string) bool {
	index := sort.SearchStrings(s.rooms, room)
	if index >= len(s.rooms) || s.rooms[index] != room {
		return false
	}
	return true
}

func (s *RoomStorage) Append(room string) (bool, error) {
	if !s.FindRoom(room) {
		return true, s.SaveAll(append(s.rooms, room))
	}
	return false, nil
}

type MessageStorage struct {
	handler Handler
}

func (s *MessageStorage) LoadAll(room string) (chat_type.Messages, error) {
	records := chat_type.Messages{}
	_, err := s.handler.Get("chatroom_"+room, &records, "chatroom")
	if err != nil {
		return records, err
	}
	messages := chat_type.Messages(make([]chat_type.Message, 0, len(records)))
	for i := range records {
		messages.Append(records[i])
	}
	return messages, err
}

func (s *MessageStorage) SaveAll(room string, messages chat_type.Messages) error {
	return s.handler.Set("chatroom_"+room, &messages, "chatroom")
}

type UserStorage struct {
	handler Handler
}

func (s *UserStorage) LoadAll() (chat_type.Users, error) {
	records := chat_type.Users{}
	_, err := s.handler.Get("user", &records)
	return records, err
}

func (s *UserStorage) SaveAll(rooms chat_type.Users) error {
	return s.handler.Set("user", rooms)
}
