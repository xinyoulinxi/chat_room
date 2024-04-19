package chat_db

import (
	chat_type "web_server/type"
)

type Handler interface {
	Set(key string, value any, group ...string) error
	Get(key string, value any, group ...string) (bool, error)
}

type RoomStorage struct {
	handler Handler
}

func (s *RoomStorage) LoadAll() (chat_type.RoomList, error) {
	records := chat_type.RoomList{}
	_, err := s.handler.Get("room_list", &records)
	return records, err
}

func (s *RoomStorage) SaveAll(rooms chat_type.RoomList) error {
	return s.handler.Set("room_list", rooms)
}

type MessageStorage struct {
	handler Handler
}

func (s *MessageStorage) LoadAll(room string) (chat_type.Messages, error) {
	records := chat_type.Messages{}
	_, err := s.handler.Get("chatroom_"+room, &records, "chatroom")
	return records, err
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
