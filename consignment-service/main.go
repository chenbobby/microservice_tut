// consignment-service/main.go
package main

import (
	"log"

	"golang.org/x/net/context"

	pb "github.com/chenbobby/microservice_tut/consignment-service/proto/consignment"
	vesselProto "github.com/chenbobby/microservice_tut/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
)

// Repository defines interface methods for `Repository`
type Repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// ConsignmentRepository - Dummy repo, this simulates a datastore of requests that a
// service is handling
type ConsignmentRepository struct {
	consignments []*pb.Consignment
}

// Create adds `consignment` to `repo.consignments`, returns `consignment` and error
func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	// TODO: remove this redundancy of slice appending
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

// GetAll returns all consignments in an IRepository
func (repo *ConsignmentRepository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// service ought to implement all the methods that satisfy the service
// defined in the protobuf definition. For the consignment service,
// we'll need the `CreateConsignment` method.
type service struct {
	repo         Repository
	vesselClient vesselProto.VesselServiceClient
}

// CreateConsignment takes a context and request and gives it to a gRPC server
// to handle
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	// Call a vessel client to find a suitable vessel for the consignment
	vesselSpec := &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	}
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), vesselSpec)
	if err != nil {
		return err
	} else {
		log.Printf("Found vessel: %s", vesselResponse.Vessel.Name)
	}

	req.VesselId = vesselResponse.Vessel.Id

	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	res.Created = true
	res.Consignment = consignment
	return nil
}

// GetConsignments takes a context and GetRequest and returns a reponse with that
// service's consignments
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &ConsignmentRepository{}

	// Create a new service through `micro`. Optionally include some options
	srv := micro.NewService(
		micro.Name("go.micro.srv.consignment"), // Must match package name given in your protobuf definition
		micro.Version("latest"),
	)

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	srv.Init()

	// Register our ShippingService with the gRPC server, tying together
	// the above code with the proto definition's auto-generated code
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})

	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
