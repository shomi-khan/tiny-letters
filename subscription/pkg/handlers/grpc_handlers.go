package handlers

import (
	"context"
	"tiny-letter/subscription/pkg/db"
	pb_subscription_manager "tiny-letter/subscription/pkg/grpc/pb/subscription_manager"
	"tiny-letter/subscription/pkg/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPC_Handler struct {
	pb_subscription_manager.UnimplementedSubscriptionManagerServer
	db *db.Repository
}

func GetGRPC_Handler(db *db.Repository) *GRPC_Handler {
	return &GRPC_Handler{
		db: db,
	}
}

func (l *GRPC_Handler) CreateSubscriptionForPublisher(c context.Context, req *pb_subscription_manager.CreateSubscriptionForPublisherRequest) (*pb_subscription_manager.CreateSubscriptionForPublisherResponse, error) {
	id, err := l.db.CreateSubscriptionForPublisher(&models.ConfirmPublisherSubscriptionRequest{
		UserId: int(req.GetUserId()),
		PlanId: int(req.GetPlanId()),
	})

	if err != nil {
		return nil, err
	}
	return &pb_subscription_manager.CreateSubscriptionForPublisherResponse{
		SubscriptionId: int32(id),
	}, nil
}

func (n *GRPC_Handler) GetContentSubscribers(c context.Context, req *pb_subscription_manager.GetContentSubscribersRequest) (*pb_subscription_manager.GetContentSubscribersResponse, error) {
	data := &models.GetContentSubscribersRequest{
		PublicationId:    int(req.PublicationId),
		ContentIsPremium: req.ContentIsPremium,
	}
	if err := data.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
	}
	subscriberIds, err := n.db.GetContentSubscribers(data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get subscriberIds: %v", err)
	}
	return &pb_subscription_manager.GetContentSubscribersResponse{
		SubscriberIds: subscriberIds,
	}, nil
}

func (n *GRPC_Handler) RollbackCreateSubscriptionForPublisher(c context.Context, req *pb_subscription_manager.RollbackCreateSubscriptionForPublisherRequest) (*pb_subscription_manager.Response, error) {
	err := n.db.RollbackCreateSubscriptionForPublisher(int(req.SubscriptionId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to rollback subscription creation: %v", err)
	}
	return &pb_subscription_manager.Response{
		IsSuccess: true,
	}, nil
}
