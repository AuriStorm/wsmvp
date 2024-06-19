package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"strconv"
	"time"

	pb "github.com/centrifugal/examples/on_demand_streams/proxyproto"
	"google.golang.org/grpc"
)

type streamerServer struct {
	pb.UnimplementedCentrifugoProxyServer
}

func (s *streamerServer) SubscribeUnidirectional(
	req *pb.SubscribeRequest,
	stream pb.CentrifugoProxy_SubscribeUnidirectionalServer,
) error {
	started := time.Now()
	fmt.Println("unidirectional subscribe called with request", req)
	defer func() {
		fmt.Println("unidirectional subscribe finished, elapsed", time.Since(started))
	}()
	stream.Send(&pb.StreamSubscribeResponse{
		SubscribeResponse: &pb.SubscribeResponse{},
	})
	i := 0
	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case <-time.After(1000 * time.Millisecond):
		}
		pub := &pb.Publication{Data: []byte(`{"input": "` + strconv.Itoa(i) + `"}`)}
		stream.Send(&pb.StreamSubscribeResponse{Publication: pub})
		i++
		if i >= 20 {
			break
		}
	}
	return nil
}

func baseClient() *http.Client {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	return client
}

var client = baseClient()

type clientData struct {
	Input string `json:"input"`
}

func (s *streamerServer) SubscribeBidirectional(
	stream pb.CentrifugoProxy_SubscribeBidirectionalServer,
) error {
	started := time.Now()
	fmt.Println("bidirectional subscribe called")
	defer func() {
		fmt.Println("bidirectional subscribe finished, elapsed", time.Since(started))
	}()
	// First message always contains SubscribeRequest.
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	fmt.Println("subscribe request received", req.SubscribeRequest)
	stream.Send(&pb.StreamSubscribeResponse{
		SubscribeResponse: &pb.SubscribeResponse{},
	})
	// The following messages contain publications from client.
	for {
		req, err := stream.Recv()
		if err != nil {
			fmt.Println(err)
			return err
		}
		data := req.Publication.Data
		fmt.Println("data from client", string(data))

		// NOTE honestly we can send payload to centrifugo directly from here, no need step to app route at all
		payload := []byte(fmt.Sprintf(`{"payload": %s}`, string(data)))
		httpreq, err := http.NewRequest(http.MethodPost, "http://app:8081/from-streams/send/", bytes.NewReader(payload))
		if err != nil {
			return err
		}
		httpreq.Header.Add("Content-Type", "application/json")
		httpresp, err := client.Do(httpreq)
		if err != nil {
			return err
		}
		if httpresp.StatusCode != http.StatusOK {
			fmt.Println(httpresp.StatusCode)
		}

		var cd clientData
		err = json.Unmarshal(data, &cd)
		if err != nil {
			return nil
		}
		pub := &pb.Publication{Data: []byte(`{"input": "` + cd.Input + `"}`)}
		stream.Send(&pb.StreamSubscribeResponse{Publication: pub})
	}
}

func main() {
	addr := ":12000"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32))
	pb.RegisterCentrifugoProxyServer(s, &streamerServer{})

	fmt.Println("Server listening on", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
