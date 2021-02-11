package services

import (
	"context"
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/bp/domain"
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	pb "gitlab.medzdrav.ru/prototype/proto/services"
	"time"
)

type deliveryServiceImpl struct {
	pb.DeliveryServiceClient
}

func newDeliveryImpl() *deliveryServiceImpl {
	a := &deliveryServiceImpl{}
	return a
}

func fromPb(d *pb.Delivery) *domain.Delivery {

	startTime := grpc.PbTSToTime(d.StartTime)

	var details map[string]interface{}
	_ = json.Unmarshal(d.Details, &details)

	return &domain.Delivery{
		Id:            d.Id,
		UserId:        d.UserId,
		ServiceTypeId: d.ServiceTypeId,
		Status:        d.Status,
		StartTime:     *startTime,
		FinishTime:    grpc.PbTSToTime(d.FinishTime),
		Details:       details,
	}
}

func (u *deliveryServiceImpl) Create(ctx context.Context, userId, serviceTypeId string, details map[string]interface{}) (*domain.Delivery, error){
	v, _ := json.Marshal(details)
	if d, err := u.DeliveryServiceClient.Create(ctx, &pb.DeliveryRequest{
		UserId:        userId,
		ServiceTypeId: serviceTypeId,
		Details:       v,
	}); err == nil {
		return fromPb(d), nil
	} else {
		return nil, err
	}
}

func (u *deliveryServiceImpl) GetDelivery(ctx context.Context, deliveryId string) (*domain.Delivery, error){
	if d, err := u.DeliveryServiceClient.GetDelivery(ctx, &pb.GetDeliveryRequest{Id: deliveryId}); err == nil {
		return fromPb(d), nil
	} else {
		return nil, err
	}
}

func (u *deliveryServiceImpl) Cancel(ctx context.Context, deliveryId string, cancelTime *time.Time) (*domain.Delivery, error){
	if d, err := u.DeliveryServiceClient.Cancel(ctx, &pb.CancelDeliveryRequest{Id: deliveryId, CancelTime: grpc.TimeToPbTS(cancelTime)}); err == nil {
		return fromPb(d), nil
	} else {
		return nil, err
	}
}

func (u *deliveryServiceImpl) Complete(ctx context.Context, deliveryId string, completeTime *time.Time) (*domain.Delivery, error){
	if d, err := u.DeliveryServiceClient.Complete(ctx, &pb.CompleteDeliveryRequest{Id: deliveryId, CompleteTime: grpc.TimeToPbTS(completeTime)}); err == nil {
		return fromPb(d), nil
	} else {
		return nil, err
	}
}

func (u *deliveryServiceImpl) UpdateDetails(ctx context.Context, id string, details map[string]interface{}) (*domain.Delivery, error) {
	v, _ := json.Marshal(details)
	if d, err := u.DeliveryServiceClient.UpdateDetails(ctx, &pb.UpdateDetailsRequest{
		Id:      id,
		Details: v,
	}); err == nil {
		return fromPb(d), nil
	} else {
		return nil, err
	}
}

