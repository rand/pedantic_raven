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
	MemoryID  string
	Content   *string
	Importance *uint32
	Tags      []string
	Context   *string
	AddTags   []string
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

	if opts.Context != nil {
		req.Context = opts.Context
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
		req.MaxResults = opts.MaxResults
	} else {
		req.MaxResults = 100 // Default limit
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
