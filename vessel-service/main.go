// vessel-service/main.go
package main

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	micro "github.com/micro/go-micro"
	pb "github.com/user/microservice_tut/vessel-service/proto/vessel"
)

type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

type VesselRepository struct {
	vessels []*pb.Vessel
}

// FindAvailable takes a sepcification and searches `repo` for a suitable
// vessel, which is then returned. Returns nil if no vessels fit specification
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range repo.vessels {
		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
			return vessel, nil
		}
	}
	errorMessage := fmt.Sprintf("No vessel found for specification (Capacity:%v\tMaxWeight:%v", spec.Capacity, spec.MaxWeight)
	return nil, errors.New(errorMessage)
}

// service needs to implement all methods of our gRPC service, which is just
// `FindAvailable`. The interface for service is applied through `Repository`.
type service struct {
	repo Repository
}

func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	// Find next available vessel
	vessel, err := s.repo.FindAvailable(req)
	if err != nil {
		return nil
	}

	// Set vessel as part of response
	res.Vessel = vessel
	return nil
}

func main() {
	vessels := []*pb.Vessel{
		&pb.Vessel{
			Id:        "vessel001",
			Name:      "Bob's Secret Vessel",
			MaxWeight: 200000,
			Capacity:  500,
		},
	}
	repo := &VesselRepository{vessels}

	srv := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)

	srv.Init()

	// Integrate protobuf service with the server
	pb.RegisterVesselServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
