package mnemosyne

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// mockMemoryServer implements pb.MemoryServiceServer for testing.
type mockMemoryServer struct {
	pb.UnimplementedMemoryServiceServer
	memories       map[string]*pb.MemoryNote
	mu             sync.RWMutex
	storeDelay     time.Duration
	recallDelay    time.Duration
	shouldFail     bool
	failCount      int
	failAfter      int
	requestCount   int
	requestCountMu sync.Mutex
}

func newMockMemoryServer() *mockMemoryServer {
	return &mockMemoryServer{
		memories: make(map[string]*pb.MemoryNote),
	}
}

func (m *mockMemoryServer) StoreMemory(ctx context.Context, req *pb.StoreMemoryRequest) (*pb.StoreMemoryResponse, error) {
	m.requestCountMu.Lock()
	m.requestCount++
	m.requestCountMu.Unlock()

	if m.storeDelay > 0 {
		select {
		case <-time.After(m.storeDelay):
		case <-ctx.Done():
			return nil, status.Error(codes.Canceled, "context canceled")
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldFail {
		return nil, status.Error(codes.Internal, "mock server error")
	}

	if m.failAfter > 0 && m.failCount < m.failAfter {
		m.failCount++
		return nil, status.Error(codes.Unavailable, "temporary failure")
	}

	// Create memory note
	memoryID := fmt.Sprintf("mem-%d", len(m.memories)+1)
	importance := uint32(5)
	if req.Importance != nil {
		importance = *req.Importance
	}

	memory := &pb.MemoryNote{
		Id:         memoryID,
		Content:    req.Content,
		Namespace:  req.Namespace,
		Importance: importance,
		Tags:       req.Tags,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	m.memories[memoryID] = memory

	return &pb.StoreMemoryResponse{Memory: memory}, nil
}

func (m *mockMemoryServer) GetMemory(ctx context.Context, req *pb.GetMemoryRequest) (*pb.GetMemoryResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.shouldFail {
		return nil, status.Error(codes.Internal, "mock server error")
	}

	memory, exists := m.memories[req.MemoryId]
	if !exists {
		return nil, status.Error(codes.NotFound, "memory not found")
	}

	return &pb.GetMemoryResponse{Memory: memory}, nil
}

func (m *mockMemoryServer) UpdateMemory(ctx context.Context, req *pb.UpdateMemoryRequest) (*pb.UpdateMemoryResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldFail {
		return nil, status.Error(codes.Internal, "mock server error")
	}

	memory, exists := m.memories[req.MemoryId]
	if !exists {
		return nil, status.Error(codes.NotFound, "memory not found")
	}

	// Update fields
	if req.Content != nil {
		memory.Content = *req.Content
	}
	if req.Importance != nil {
		memory.Importance = *req.Importance
	}
	if len(req.Tags) > 0 {
		memory.Tags = req.Tags
	}
	if len(req.AddTags) > 0 {
		memory.Tags = append(memory.Tags, req.AddTags...)
	}
	memory.UpdatedAt = time.Now().Unix()

	return &pb.UpdateMemoryResponse{Memory: memory}, nil
}

func (m *mockMemoryServer) DeleteMemory(ctx context.Context, req *pb.DeleteMemoryRequest) (*pb.DeleteMemoryResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldFail {
		return nil, status.Error(codes.Internal, "mock server error")
	}

	if _, exists := m.memories[req.MemoryId]; !exists {
		return nil, status.Error(codes.NotFound, "memory not found")
	}

	delete(m.memories, req.MemoryId)
	return &pb.DeleteMemoryResponse{}, nil
}

func (m *mockMemoryServer) ListMemories(ctx context.Context, req *pb.ListMemoriesRequest) (*pb.ListMemoriesResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.shouldFail {
		return nil, status.Error(codes.Internal, "mock server error")
	}

	var results []*pb.MemoryNote
	for _, memory := range m.memories {
		// Filter by namespace
		if req.Namespace != nil && !namespacesMatch(memory.Namespace, req.Namespace) {
			continue
		}

		// Filter by min importance
		if req.MinImportance != nil && memory.Importance < *req.MinImportance {
			continue
		}

		// Filter by tags
		if len(req.Tags) > 0 && !hasAllTags(memory.Tags, req.Tags) {
			continue
		}

		results = append(results, memory)
	}

	// Apply limit
	if req.Limit > 0 && uint32(len(results)) > req.Limit {
		results = results[:req.Limit]
	}

	return &pb.ListMemoriesResponse{Memories: results}, nil
}

func (m *mockMemoryServer) Recall(ctx context.Context, req *pb.RecallRequest) (*pb.RecallResponse, error) {
	m.requestCountMu.Lock()
	m.requestCount++
	m.requestCountMu.Unlock()

	if m.recallDelay > 0 {
		select {
		case <-time.After(m.recallDelay):
		case <-ctx.Done():
			return nil, status.Error(codes.Canceled, "context canceled")
		}
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.shouldFail {
		return nil, status.Error(codes.Internal, "mock server error")
	}

	var results []*pb.SearchResult
	for _, memory := range m.memories {
		// Filter by namespace
		if req.Namespace != nil && !namespacesMatch(memory.Namespace, req.Namespace) {
			continue
		}

		// Filter by min importance
		if req.MinImportance != nil && memory.Importance < *req.MinImportance {
			continue
		}

		// Filter by tags
		if len(req.Tags) > 0 && !hasAllTags(memory.Tags, req.Tags) {
			continue
		}

		// Create search result with mock score
		result := &pb.SearchResult{
			Memory:         memory,
			RelevanceScore: 0.85,
			SemanticScore:  0.9,
			FtsScore:       0.8,
			GraphScore:     0.7,
		}
		results = append(results, result)
	}

	// Apply limit
	maxResults := req.MaxResults
	if maxResults == 0 {
		maxResults = 10
	}
	if uint32(len(results)) > maxResults {
		results = results[:maxResults]
	}

	return &pb.RecallResponse{Results: results}, nil
}

func (m *mockMemoryServer) SemanticSearch(ctx context.Context, req *pb.SemanticSearchRequest) (*pb.SemanticSearchResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.shouldFail {
		return nil, status.Error(codes.Internal, "mock server error")
	}

	var results []*pb.SearchResult
	for _, memory := range m.memories {
		// Filter by namespace
		if req.Namespace != nil && !namespacesMatch(memory.Namespace, req.Namespace) {
			continue
		}

		result := &pb.SearchResult{
			Memory:         memory,
			RelevanceScore: 0.9,
			SemanticScore:  0.95,
		}
		results = append(results, result)
	}

	// Apply limit
	maxResults := req.MaxResults
	if maxResults == 0 {
		maxResults = 10
	}
	if uint32(len(results)) > maxResults {
		results = results[:maxResults]
	}

	return &pb.SemanticSearchResponse{Results: results}, nil
}

func (m *mockMemoryServer) GraphTraverse(ctx context.Context, req *pb.GraphTraverseRequest) (*pb.GraphTraverseResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.shouldFail {
		return nil, status.Error(codes.Internal, "mock server error")
	}

	var memories []*pb.MemoryNote
	var edges []*pb.GraphEdge

	// Return seed memories
	for _, seedID := range req.SeedIds {
		if memory, exists := m.memories[seedID]; exists {
			memories = append(memories, memory)
		}
	}

	return &pb.GraphTraverseResponse{Memories: memories, Edges: edges}, nil
}

func (m *mockMemoryServer) GetContext(ctx context.Context, req *pb.GetContextRequest) (*pb.GetContextResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.shouldFail {
		return nil, status.Error(codes.Internal, "mock server error")
	}

	var memories []*pb.MemoryNote
	for _, memID := range req.MemoryIds {
		if memory, exists := m.memories[memID]; exists {
			memories = append(memories, memory)
		}
	}

	return &pb.GetContextResponse{Memories: memories}, nil
}

// mockHealthServer implements pb.HealthServiceServer for testing.
type mockHealthServer struct {
	pb.UnimplementedHealthServiceServer
	healthy   bool
	mu        sync.RWMutex
	checkDelay time.Duration
}

func newMockHealthServer() *mockHealthServer {
	return &mockHealthServer{
		healthy: true,
	}
}

func (m *mockHealthServer) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	if m.checkDelay > 0 {
		select {
		case <-time.After(m.checkDelay):
		case <-ctx.Done():
			return nil, status.Error(codes.Canceled, "context canceled")
		}
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.healthy {
		return nil, status.Error(codes.Unavailable, "service unhealthy")
	}

	return &pb.HealthCheckResponse{
		Status: pb.HealthStatus_SERVING,
	}, nil
}

func (m *mockHealthServer) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.healthy {
		return nil, status.Error(codes.Unavailable, "service unhealthy")
	}

	stats := &pb.Stats{
		TotalMemories:  100,
		TotalNamespaces: 5,
		DatabaseSizeMb: 50.5,
	}

	return &pb.GetStatsResponse{Stats: stats}, nil
}

// testServer combines both services for testing.
type testServer struct {
	grpcServer *grpc.Server
	listener   net.Listener
	address    string
	memory     *mockMemoryServer
	health     *mockHealthServer
}

func newTestServer() (*testServer, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	memoryServer := newMockMemoryServer()
	healthServer := newMockHealthServer()

	pb.RegisterMemoryServiceServer(grpcServer, memoryServer)
	pb.RegisterHealthServiceServer(grpcServer, healthServer)

	ts := &testServer{
		grpcServer: grpcServer,
		listener:   listener,
		address:    listener.Addr().String(),
		memory:     memoryServer,
		health:     healthServer,
	}

	go grpcServer.Serve(listener)

	return ts, nil
}

func (ts *testServer) Stop() {
	if ts.grpcServer != nil {
		ts.grpcServer.Stop()
	}
	if ts.listener != nil {
		ts.listener.Close()
	}
}

// Helper functions

func namespacesMatch(ns1, ns2 *pb.Namespace) bool {
	if ns1 == nil || ns2 == nil {
		return ns1 == ns2
	}

	switch n1 := ns1.Namespace.(type) {
	case *pb.Namespace_Global:
		_, ok := ns2.Namespace.(*pb.Namespace_Global)
		return ok
	case *pb.Namespace_Project:
		n2, ok := ns2.Namespace.(*pb.Namespace_Project)
		return ok && n1.Project.Name == n2.Project.Name
	case *pb.Namespace_Session:
		n2, ok := ns2.Namespace.(*pb.Namespace_Session)
		return ok && n1.Session.Project == n2.Session.Project && n1.Session.SessionId == n2.Session.SessionId
	}
	return false
}

func hasAllTags(memoryTags, requiredTags []string) bool {
	tagSet := make(map[string]bool)
	for _, tag := range memoryTags {
		tagSet[tag] = true
	}

	for _, required := range requiredTags {
		if !tagSet[required] {
			return false
		}
	}
	return true
}
