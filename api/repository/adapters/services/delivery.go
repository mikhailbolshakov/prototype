package services

import (
	"context"
	"encoding/json"
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

func (u *deliveryServiceImpl) Create(ctx context.Context, userId, serviceTypeId string, details map[string]interface{}) (*pb.Delivery, error){
	v, _ := json.Marshal(details)
	return u.DeliveryServiceClient.Create(ctx, &pb.DeliveryRequest{
		UserId:        userId,
		ServiceTypeId: serviceTypeId,
		Details:       v,
	})
}

func (u *deliveryServiceImpl) GetDelivery(ctx context.Context, deliveryId string) (*pb.Delivery, error){
	return u.DeliveryServiceClient.GetDelivery(ctx, &pb.GetDeliveryRequest{Id: deliveryId})
}

func (u *deliveryServiceImpl) Cancel(ctx context.Context, deliveryId string, cancelTime *time.Time) (*pb.Delivery, error){
	return u.DeliveryServiceClient.Cancel(ctx, &pb.CancelDeliveryRequest{Id: deliveryId, CancelTime: grpc.TimeToPbTS(cancelTime)})
}

func (u *deliveryServiceImpl) Complete(ctx context.Context, deliveryId string, completeTime *time.Time) (*pb.Delivery, error){
	return u.DeliveryServiceClient.Complete(ctx, &pb.CompleteDeliveryRequest{Id: deliveryId, CompleteTime: grpc.TimeToPbTS(completeTime)})
}

func (u *deliveryServiceImpl) UpdateDetails(ctx context.Context, id string, details map[string]interface{}) (*pb.Delivery, error) {
	v, _ := json.Marshal(details)
	return u.DeliveryServiceClient.UpdateDetails(ctx, &pb.UpdateDetailsRequest{
		Id:      id,
		Details: v,
	})
}

