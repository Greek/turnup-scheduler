package server

import (
	"context"
	"encoding/json"
	"fmt"
	"turnup-scheduler/internal/lib"
	"turnup-scheduler/internal/logging"
	"turnup-scheduler/pkg/scheduler"
	pb "turnup-scheduler/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetSnapshot(_ context.Context, in *pb.GetSnapshotRequest) (*pb.GetSnapshotResponse, error) {
	log := logging.BuildLogger("GetSnapshot")
	currDate := lib.BuildDate()
	namespace := in.GetNamespace()

	if namespace == "" {
		return nil, status.Error(codes.InvalidArgument, "namespace cannot be nil")
	}

	log.Info("Received namespace: " + in.GetNamespace())
	createdSnapshot, err := s.Scheduler.CreateSnapshot(currDate, namespace, scheduler.CreateSnapshotOpts{Overwrite: in.GetOverwrite()})
	if err != nil {
		// If snapshot already exists, we fetch it and return it.
		if err.Error() == "key already exists" {
			existingSnapshot, err2 := s.Scheduler.GetSnapshot(currDate, namespace)
			if err2 != nil {
				fmt.Print(err2)
				return nil, err2
			}

			bytes, marshalErr := json.Marshal(existingSnapshot)
			if marshalErr != nil {
				fmt.Print(marshalErr)
				return nil, marshalErr
			}

			return &pb.GetSnapshotResponse{
				Snapshot: string(bytes),
			}, nil
		}

	}
	return &pb.GetSnapshotResponse{
		Snapshot: createdSnapshot,
	}, nil
}
