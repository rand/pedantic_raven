package mnemosyne

import (
	"context"
	"fmt"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// StoreMemoryOptions holds options for storing a memory.
type StoreMemoryOptions struct {
	Content          string
	Namespace        *pb.Namespace
	Importance       *uint32 // Optional: 1-10
	Context          string
	Tags             []string
	MemoryType       *pb.MemoryType
	SkipLLMEnrichment bool
}

// StoreMemory stores a new memory note.
func (c *Client) StoreMemory(ctx context.Context, opts StoreMemoryOptions) (*pb.MemoryNote, error) {
	if !c.connected {
		return nil, ErrNotConnected
	}

	if opts.Content == "" {
		return nil, fmt.Errorf("%w: content is required", ErrInvalidArgument)
	}

	if opts.Namespace == nil {
		return nil, fmt.Errorf("%w: namespace is required", ErrInvalidArgument)
	}

	req := &pb.StoreMemoryRequest{
		Content:           opts.Content,
		Namespace:         opts.Namespace,
		Context:           &opts.Context,
		Tags:              opts.Tags,
		SkipLlmEnrichment: opts.SkipLLMEnrichment,
	}

	if opts.Importance != nil {
		req.Importance = opts.Importance
	}

	if opts.MemoryType != nil {
		req.MemoryType = opts.MemoryType
	}

	resp, err := c.memoryClient.StoreMemory(ctx, req)
	if err != nil {
		return nil, wrapError(err, "store memory")
	}

	return resp.Memory, nil
}

// GetMemory retrieves a memory by ID.
func (c *Client) GetMemory(ctx context.Context, memoryID string) (*pb.MemoryNote, error) {
	if !c.connected {
		return nil, ErrNotConnected
	}

	if memoryID == "" {
		return nil, fmt.Errorf("%w: memory ID is required", ErrInvalidArgument)
	}

	req := &pb.GetMemoryRequest{
		MemoryId: memoryID,
	}

	resp, err := c.memoryClient.GetMemory(ctx, req)
	if err != nil {
		return nil, wrapError(err, "get memory")
	}

	return resp.Memory, nil
}

// UpdateMemoryOptions holds options for updating a memory.
type UpdateMemoryOptions struct {
	MemoryID   string
	Content    *string
	Importance *uint32
	Tags       []string
	AddTags    []string
	RemoveTags []string
}

// UpdateMemory updates an existing memory note.
func (c *Client) UpdateMemory(ctx context.Context, opts UpdateMemoryOptions) (*pb.MemoryNote, error) {
	if !c.connected {
		return nil, ErrNotConnected
	}

	if opts.MemoryID == "" {
		return nil, fmt.Errorf("%w: memory ID is required", ErrInvalidArgument)
	}

	req := &pb.UpdateMemoryRequest{
		MemoryId: opts.MemoryID,
	}

	if opts.Content != nil {
		req.Content = opts.Content
	}

	if opts.Importance != nil {
		req.Importance = opts.Importance
	}

	if len(opts.Tags) > 0 {
		req.Tags = opts.Tags
	}

	if len(opts.AddTags) > 0 {
		req.AddTags = opts.AddTags
	}

	if len(opts.RemoveTags) > 0 {
		req.RemoveTags = opts.RemoveTags
	}

	resp, err := c.memoryClient.UpdateMemory(ctx, req)
	if err != nil {
		return nil, wrapError(err, "update memory")
	}

	return resp.Memory, nil
}

// DeleteMemory deletes a memory by ID.
func (c *Client) DeleteMemory(ctx context.Context, memoryID string) error {
	if !c.connected {
		return ErrNotConnected
	}

	if memoryID == "" {
		return fmt.Errorf("%w: memory ID is required", ErrInvalidArgument)
	}

	req := &pb.DeleteMemoryRequest{
		MemoryId: memoryID,
	}

	_, err := c.memoryClient.DeleteMemory(ctx, req)
	if err != nil {
		return wrapError(err, "delete memory")
	}

	return nil
}

// ListMemoriesOptions holds options for listing memories.
type ListMemoriesOptions struct {
	Namespace       *pb.Namespace
	MemoryTypes     []pb.MemoryType
	Tags            []string
	MinImportance   *uint32
	MaxResults      uint32
	IncludeArchived bool
}

// ListMemories retrieves a list of memories with optional filtering.
func (c *Client) ListMemories(ctx context.Context, opts ListMemoriesOptions) ([]*pb.MemoryNote, error) {
	if !c.connected {
		return nil, ErrNotConnected
	}

	req := &pb.ListMemoriesRequest{
		Namespace:       opts.Namespace,
		MemoryTypes:     opts.MemoryTypes,
		Tags:            opts.Tags,
		IncludeArchived: opts.IncludeArchived,
	}

	if opts.MinImportance != nil {
		req.MinImportance = opts.MinImportance
	}

	if opts.MaxResults > 0 {
		req.Limit = opts.MaxResults
	} else {
		req.Limit = 100 // Default limit
	}

	resp, err := c.memoryClient.ListMemories(ctx, req)
	if err != nil {
		return nil, wrapError(err, "list memories")
	}

	return resp.Memories, nil
}

// Helper functions for creating common namespace types

// GlobalNamespace creates a global namespace.
func GlobalNamespace() *pb.Namespace {
	return &pb.Namespace{
		Namespace: &pb.Namespace_Global{
			Global: &pb.GlobalNamespace{},
		},
	}
}

// ProjectNamespace creates a project namespace.
func ProjectNamespace(name string) *pb.Namespace {
	return &pb.Namespace{
		Namespace: &pb.Namespace_Project{
			Project: &pb.ProjectNamespace{
				Name: name,
			},
		},
	}
}

// SessionNamespace creates a session namespace.
func SessionNamespace(project, sessionID string) *pb.Namespace {
	return &pb.Namespace{
		Namespace: &pb.Namespace_Session{
			Session: &pb.SessionNamespace{
				Project:   project,
				SessionId: sessionID,
			},
		},
	}
}

// RecallOptions holds options for hybrid search.
type RecallOptions struct {
	Query           string
	Namespace       *pb.Namespace
	MaxResults      uint32
	MinImportance   *uint32
	MemoryTypes     []pb.MemoryType
	Tags            []string
	IncludeArchived bool
	SemanticWeight  *float32 // Default: 0.7
	FtsWeight       *float32 // Default: 0.2
	GraphWeight     *float32 // Default: 0.1
}

// Recall performs hybrid search (semantic + FTS + graph).
func (c *Client) Recall(ctx context.Context, opts RecallOptions) ([]*pb.SearchResult, error) {
	if !c.connected {
		return nil, ErrNotConnected
	}

	if opts.Query == "" {
		return nil, fmt.Errorf("%w: query is required", ErrInvalidArgument)
	}

	req := &pb.RecallRequest{
		Query:           opts.Query,
		Namespace:       opts.Namespace,
		MaxResults:      opts.MaxResults,
		MemoryTypes:     opts.MemoryTypes,
		Tags:            opts.Tags,
		IncludeArchived: opts.IncludeArchived,
	}

	if opts.MinImportance != nil {
		req.MinImportance = opts.MinImportance
	}

	if opts.SemanticWeight != nil {
		req.SemanticWeight = opts.SemanticWeight
	}

	if opts.FtsWeight != nil {
		req.FtsWeight = opts.FtsWeight
	}

	if opts.GraphWeight != nil {
		req.GraphWeight = opts.GraphWeight
	}

	if opts.MaxResults == 0 {
		req.MaxResults = 10 // Default limit
	}

	resp, err := c.memoryClient.Recall(ctx, req)
	if err != nil {
		return nil, wrapError(err, "recall")
	}

	return resp.Results, nil
}

// SemanticSearchOptions holds options for semantic search.
type SemanticSearchOptions struct {
	Embedding       []float32
	Namespace       *pb.Namespace
	MaxResults      uint32
	MinImportance   *uint32
	IncludeArchived bool
}

// SemanticSearch performs pure semantic search using an embedding vector.
func (c *Client) SemanticSearch(ctx context.Context, opts SemanticSearchOptions) ([]*pb.SearchResult, error) {
	if !c.connected {
		return nil, ErrNotConnected
	}

	if len(opts.Embedding) == 0 {
		return nil, fmt.Errorf("%w: embedding is required", ErrInvalidArgument)
	}

	req := &pb.SemanticSearchRequest{
		Embedding:       opts.Embedding,
		Namespace:       opts.Namespace,
		MaxResults:      opts.MaxResults,
		IncludeArchived: opts.IncludeArchived,
	}

	if opts.MinImportance != nil {
		req.MinImportance = opts.MinImportance
	}

	if opts.MaxResults == 0 {
		req.MaxResults = 10 // Default limit
	}

	resp, err := c.memoryClient.SemanticSearch(ctx, req)
	if err != nil {
		return nil, wrapError(err, "semantic search")
	}

	return resp.Results, nil
}

// GraphTraverseOptions holds options for graph traversal.
type GraphTraverseOptions struct {
	SeedIDs         []string
	MaxHops         uint32
	LinkTypes       []pb.LinkType
	MinLinkStrength *float32
}

// GraphTraverse traverses the memory graph from seed nodes.
func (c *Client) GraphTraverse(ctx context.Context, opts GraphTraverseOptions) ([]*pb.MemoryNote, []*pb.GraphEdge, error) {
	if !c.connected {
		return nil, nil, ErrNotConnected
	}

	if len(opts.SeedIDs) == 0 {
		return nil, nil, fmt.Errorf("%w: at least one seed ID is required", ErrInvalidArgument)
	}

	req := &pb.GraphTraverseRequest{
		SeedIds:   opts.SeedIDs,
		MaxHops:   opts.MaxHops,
		LinkTypes: opts.LinkTypes,
	}

	if opts.MinLinkStrength != nil {
		req.MinLinkStrength = opts.MinLinkStrength
	}

	if opts.MaxHops == 0 {
		req.MaxHops = 2 // Default max hops
	}

	resp, err := c.memoryClient.GraphTraverse(ctx, req)
	if err != nil {
		return nil, nil, wrapError(err, "graph traverse")
	}

	return resp.Memories, resp.Edges, nil
}

// GetContextOptions holds options for getting memory context.
type GetContextOptions struct {
	MemoryIDs       []string
	IncludeLinks    bool
	MaxLinkedDepth  uint32
}

// GetContext retrieves memories with their surrounding context.
func (c *Client) GetContext(ctx context.Context, opts GetContextOptions) (*pb.GetContextResponse, error) {
	if !c.connected {
		return nil, ErrNotConnected
	}

	if len(opts.MemoryIDs) == 0 {
		return nil, fmt.Errorf("%w: at least one memory ID is required", ErrInvalidArgument)
	}

	req := &pb.GetContextRequest{
		MemoryIds:      opts.MemoryIDs,
		IncludeLinks:   opts.IncludeLinks,
		MaxLinkedDepth: opts.MaxLinkedDepth,
	}

	if opts.MaxLinkedDepth == 0 {
		req.MaxLinkedDepth = 1 // Default depth
	}

	resp, err := c.memoryClient.GetContext(ctx, req)
	if err != nil {
		return nil, wrapError(err, "get context")
	}

	return resp, nil
}

// RecallStream performs streaming hybrid search.
func (c *Client) RecallStream(ctx context.Context, opts RecallOptions) (pb.MemoryService_RecallStreamClient, error) {
	if !c.connected {
		return nil, ErrNotConnected
	}

	if opts.Query == "" {
		return nil, fmt.Errorf("%w: query is required", ErrInvalidArgument)
	}

	req := &pb.RecallRequest{
		Query:           opts.Query,
		Namespace:       opts.Namespace,
		MaxResults:      opts.MaxResults,
		MemoryTypes:     opts.MemoryTypes,
		Tags:            opts.Tags,
		IncludeArchived: opts.IncludeArchived,
	}

	if opts.MinImportance != nil {
		req.MinImportance = opts.MinImportance
	}

	if opts.SemanticWeight != nil {
		req.SemanticWeight = opts.SemanticWeight
	}

	if opts.FtsWeight != nil {
		req.FtsWeight = opts.FtsWeight
	}

	if opts.GraphWeight != nil {
		req.GraphWeight = opts.GraphWeight
	}

	if opts.MaxResults == 0 {
		req.MaxResults = 10
	}

	stream, err := c.memoryClient.RecallStream(ctx, req)
	if err != nil {
		return nil, wrapError(err, "recall stream")
	}

	return stream, nil
}

// ListMemoriesStream performs streaming list of memories.
func (c *Client) ListMemoriesStream(ctx context.Context, opts ListMemoriesOptions) (pb.MemoryService_ListMemoriesStreamClient, error) {
	if !c.connected {
		return nil, ErrNotConnected
	}

	req := &pb.ListMemoriesRequest{
		Namespace:       opts.Namespace,
		MemoryTypes:     opts.MemoryTypes,
		Tags:            opts.Tags,
		IncludeArchived: opts.IncludeArchived,
	}

	if opts.MinImportance != nil {
		req.MinImportance = opts.MinImportance
	}

	if opts.MaxResults > 0 {
		req.Limit = opts.MaxResults
	} else {
		req.Limit = 100 // Default limit
	}

	stream, err := c.memoryClient.ListMemoriesStream(ctx, req)
	if err != nil {
		return nil, wrapError(err, "list memories stream")
	}

	return stream, nil
}

// StoreMemoryStream stores a memory with progress updates.
func (c *Client) StoreMemoryStream(ctx context.Context, opts StoreMemoryOptions) (pb.MemoryService_StoreMemoryStreamClient, error) {
	if !c.connected {
		return nil, ErrNotConnected
	}

	if opts.Content == "" {
		return nil, fmt.Errorf("%w: content is required", ErrInvalidArgument)
	}

	if opts.Namespace == nil {
		return nil, fmt.Errorf("%w: namespace is required", ErrInvalidArgument)
	}

	req := &pb.StoreMemoryRequest{
		Content:           opts.Content,
		Namespace:         opts.Namespace,
		Context:           &opts.Context,
		Tags:              opts.Tags,
		SkipLlmEnrichment: opts.SkipLLMEnrichment,
	}

	if opts.Importance != nil {
		req.Importance = opts.Importance
	}

	if opts.MemoryType != nil {
		req.MemoryType = opts.MemoryType
	}

	stream, err := c.memoryClient.StoreMemoryStream(ctx, req)
	if err != nil {
		return nil, wrapError(err, "store memory stream")
	}

	return stream, nil
}
