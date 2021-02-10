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

func (u *deliveryServiceImpl) Create(userId, serviceTypeId string, details map[string]interface{}) (*pb.Delivery, error){
	v, _ := json.Marshal(details)
	return u.DeliveryServiceClient.Create(context.Background(), &pb.DeliveryRequest{
		UserId:        userId,
		ServiceTypeId: serviceTypeId,
		Details:       v,
	})
}

func (u *deliveryServiceImpl) GetDelivery(deliveryId string) (*pb.Delivery, error){
	return u.DeliveryServiceClient.GetDelivery(context.Background(), &pb.GetDeliveryRequest{Id: deliveryId})
}

func (u *deliveryServiceImpl) Cancel(deliveryId string, cancelTime *time.Time) (*pb.Delivery, error){
	return u.DeliveryServiceClient.Cancel(context.Background(), &pb.CancelDeliveryRequest{Id: deliveryId, CancelTime: grpc.TimeToPbTS(cancelTime)})
}

func (u *deliveryServiceImpl) Complete(deliveryId string, completeTime *time.Time) (*pb.Delivery, error){
	return u.DeliveryServiceClient.Complete(context.Background(), &pb.CompleteDeliveryRequest{Id: deliveryId, CompleteTime: grpc.TimeToPbTS(completeTime)})
}

func (u *deliveryServiceImpl) UpdateDetails(id string, details map[string]interface{}) (*pb.Delivery, error) {
	v, _ := json.Marshal(details)
	return u.DeliveryServiceClient.UpdateDetails(context.Background(), &pb.UpdateDetailsRequest{
		Id:      id,
		Details: v,
	})
}

