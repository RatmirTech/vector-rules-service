package grpc

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/ratmirtech/vector-rules-service/internal/domain"
	"github.com/ratmirtech/vector-rules-service/internal/transport/grpc/pb"
)

// ruleRetrievalServer implements the gRPC RuleRetrieval service
type ruleRetrievalServer struct {
	pb.UnimplementedRuleRetrievalServiceServer
	ruleService domain.RuleService
}

// NewRuleRetrievalServer creates a new gRPC server for rule retrieval
func NewRuleRetrievalServer(ruleService domain.RuleService) pb.RuleRetrievalServiceServer {
	return &ruleRetrievalServer{
		ruleService: ruleService,
	}
}

// Retrieve implements the gRPC Retrieve method for vector similarity search
func (s *ruleRetrievalServer) Retrieve(ctx context.Context, req *pb.RetrieveRequest) (*pb.RetrieveResponse, error) {
	// Validate request
	if req.N <= 0 {
		return nil, fmt.Errorf("n must be greater than 0")
	}
	
	if len(req.Queries) == 0 {
		return nil, fmt.Errorf("queries cannot be empty")
	}

	// Convert to domain query
	query := &domain.RetrieveRulesQuery{
		N:       int(req.N),
		Queries: req.Queries,
	}
	
	if req.Type != nil {
		query.Type = req.Type
	}

	// Call business logic
	matches, err := s.ruleService.RetrieveSimilar(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve similar rules: %w", err)
	}

	// Convert domain matches to protobuf response
	response := &pb.RetrieveResponse{
		Rules: make([]*pb.RuleMatch, len(matches)),
	}

	for i, match := range matches {
		// Convert JSON content to protobuf Struct
		var contentMap map[string]interface{}
		if err := json.Unmarshal(match.Content, &contentMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal rule content: %w", err)
		}

		contentStruct, err := structpb.NewStruct(contentMap)
		if err != nil {
			return nil, fmt.Errorf("failed to create protobuf struct: %w", err)
		}

		ruleTypeName := ""
		if match.RuleTypeName != nil {
			ruleTypeName = *match.RuleTypeName
		}

		response.Rules[i] = &pb.RuleMatch{
			Id:        match.ID,
			Type:      ruleTypeName,
			Content:   contentStruct,
			Score:     match.Score,
			CreatedAt: match.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: match.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return response, nil
}