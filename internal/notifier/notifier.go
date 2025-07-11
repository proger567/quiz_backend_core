package notifier

import (
	"context"
	proto "github.com/proger567/quiz_protos/gen/go/notifier"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"quiz_backend_core/internal/model"
)

// Notifier TODO check interface implementation
type Notifier struct {
	logger *logrus.Logger
	conn   *grpc.ClientConn
	client proto.NotifierClient
}

func NewNotifier(target string, logger *logrus.Logger) (model.Notifier, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := proto.NewNotifierClient(conn)

	return &Notifier{
		logger: logger,
		conn:   conn,
		client: client,
	}, nil
}

// Notify TODO check interface implementation
func (n *Notifier) Notify(ctx context.Context, room, event, data string) error {
	//if you will need to add metadata, your can use next code

	//ctx, err := fillGRPCMetadata(ctx)
	//if err != nil {
	//	return err
	//}

	_, err := n.client.Notify(ctx, &proto.NotifyRequest{
		Room:  room,
		Event: event,
		Data:  data,
	})

	return err
}

func (n *Notifier) Close() error {
	return n.conn.Close()
}

//func fillGRPCMetadata(ctx context.Context) (context.Context, error) {
//	userID, ok1 := ctx.Value(constants.ContextVariablesUserID).(int64)
//	userRole, ok2 := ctx.Value(constants.ContextVariablesUserRole).(dto.Role) // TODO ALWAYS ok (add check!)
//	if !ok1 || !ok2 {
//		return nil, errors.New("bad user id or role") //TODO
//	}
//
//	md := metadata.New(map[string]string{
//		"X-User-ID":   fmt.Sprintf("%d", userID),
//		"X-User-Role": string(userRole),
//	})
//	return metadata.NewOutgoingContext(ctx, md), nil
//}
