package proto

import (
	"context"
	"encoding/json"

	"github.com/oliver7100/advertisement-service/database"
)

type service struct {
	UnimplementedAdvertisementServiceServer
	Conn *database.Connection
}

func TypeConverter[R any](data any) (*R, error) {
	var result R
	b, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}

func (s *service) CreateAdvertisement(ctx context.Context, req *Advertisement) (*Advertisement, error) {
	m, err := TypeConverter[database.Advertisement](req)

	if err != nil {
		return nil, err
	}

	if tx := s.Conn.Instance.Create(&m); tx.Error != nil {
		return nil, tx.Error
	}

	return req, nil
}

func (s *service) GetAdvertisements(ctx context.Context, req *GetAllAdvertisementsRequest) (*GetAllAdvertisementsResponse, error) {
	arr := new([]*Advertisement)

	if tx := s.Conn.Instance.Model(&database.Advertisement{}).Find(arr, "activated = ?", req.AllowDisabled); tx.Error != nil {
		return nil, tx.Error
	}

	return &GetAllAdvertisementsResponse{
		Items: *arr,
	}, nil
}

func NewService(conn *database.Connection) *service {
	return &service{
		Conn: conn,
	}
}
