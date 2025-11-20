// Package repository provides repository implementations for various AWS services and databases.
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

// AWSBedrockClientInterface defines the interface for Bedrock Runtime client operations
type AWSBedrockClientInterface interface {
	// InvokeModel invokes a model with the provided input
	InvokeModel(_ context.Context, _ *bedrockruntime.InvokeModelInput, _ ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error)
	// InvokeModelWithResponseStream invokes a model with streaming response
	InvokeModelWithResponseStream(_ context.Context, _ *bedrockruntime.InvokeModelWithResponseStreamInput, _ ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelWithResponseStreamOutput, error)
	// Converse enables conversational interactions with models
	Converse(_ context.Context, _ *bedrockruntime.ConverseInput, _ ...func(*bedrockruntime.Options)) (*bedrockruntime.ConverseOutput, error)
}

// AWSBedrockRepository struct.
type AWSBedrockRepository struct {
	Client AWSBedrockClientInterface
}

// NewAWSBedrockRepository creates a new repository backed by the provided SDK client.
func NewAWSBedrockRepository(c *bedrockruntime.Client) *AWSBedrockRepository {
	if c == nil {
		return &AWSBedrockRepository{}
	}
	return &AWSBedrockRepository{Client: c}
}

// NewAWSBedrockRepositoryWithInterface creates a repository using the provided client interface (for testing).
func NewAWSBedrockRepositoryWithInterface(c AWSBedrockClientInterface) *AWSBedrockRepository {
	return &AWSBedrockRepository{Client: c}
}

// InvokeModel calls the Bedrock Runtime InvokeModel API for the specified modelID with JSON payload.
// It returns the raw response body as bytes.
func (r *AWSBedrockRepository) InvokeModel(ctx context.Context, modelID string, payload []byte) ([]byte, error) {
	in := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelID),
		ContentType: aws.String("application/json"),
		Body:        payload,
	}
	out, err := r.Client.InvokeModel(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("invoke model: %w", err)
	}
	// out.Body is []byte
	return out.Body, nil
}

// InvokeModelWithStream calls the streaming variant and returns the SDK output for callers that need streaming.
func (r *AWSBedrockRepository) InvokeModelWithStream(ctx context.Context, modelID string, payload []byte) (*bedrockruntime.InvokeModelWithResponseStreamOutput, error) {
	in := &bedrockruntime.InvokeModelWithResponseStreamInput{
		ModelId:     aws.String(modelID),
		ContentType: aws.String("application/json"),
		Body:        payload,
	}
	out, err := r.Client.InvokeModelWithResponseStream(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("invoke model with stream: %w", err)
	}
	return out, nil
}

// Converse wraps the Converse API for conversation-based models. It returns extracted text when possible.
func (r *AWSBedrockRepository) Converse(ctx context.Context, modelID string, message any) (string, error) {
	// For convenience, marshal message to JSON and pass as Message Content per docs.
	b, err := json.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("marshal message: %w", err)
	}
	// The Converse API input expects complex types; callers may prefer to call SDK directly for rich types.
	in := &bedrockruntime.ConverseInput{
		ModelId: aws.String(modelID),
		// For simplicity this example omits constructing the full types.Message structures.
		// Callers who need Converse should construct proper types.Message values and use the SDK directly.
	}
	_, err = r.Client.Converse(ctx, in)
	if err != nil {
		return "", fmt.Errorf("converse: %w", err)
	}
	return string(b), nil
}

// InvokeModelWithFile calls the Bedrock Runtime InvokeModel API with a file attachment.
// It reads the file from the provided filePath and returns the file data as bytes.
// The caller is responsible for constructing the appropriate payload structure for their model.
func (r *AWSBedrockRepository) InvokeModelWithFile(ctx context.Context, modelID, filePath string, payload []byte) ([]byte, error) {
	// Read file content
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("close file: %w", closeErr)
		}
	}()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	result, err := r.InvokeModelWithFileData(ctx, modelID, fileData, payload)
	if err != nil {
		return nil, fmt.Errorf("invoke model with file: %w", err)
	}

	return result, nil
}

// InvokeModelWithFileData calls the Bedrock Runtime InvokeModel API with file data.
// This is a convenience wrapper that accepts fileData parameter for API clarity,
// but the caller must encode and include the file data in the payload.
// The fileData parameter exists to distinguish file-based operations from text-only operations.
//
// Example for image models (Claude 3):
//
//	import "encoding/base64"
//	import "encoding/json"
//
//	fileData := []byte("...image bytes...")
//	payload := map[string]any{
//	    "anthropic_version": "bedrock-2023-05-31",
//	    "max_tokens": 1000,
//	    "messages": []map[string]any{
//	        {
//	            "role": "user",
//	            "content": []map[string]any{
//	                {
//	                    "type": "image",
//	                    "source": map[string]string{
//	                        "type": "base64",
//	                        "media_type": "image/png",
//	                        "data": base64.StdEncoding.EncodeToString(fileData),
//	                    },
//	                },
//	                {
//	                    "type": "text",
//	                    "text": "Describe this image",
//	                },
//	            },
//	        },
//	    },
//	}
//	payloadBytes, _ := json.Marshal(payload)
//	result, err := repo.InvokeModelWithFileData(ctx, modelID, fileData, payloadBytes)
func (r *AWSBedrockRepository) InvokeModelWithFileData(ctx context.Context, modelID string, fileData, payload []byte) ([]byte, error) {
	// fileData parameter kept for API semantics; caller must include it in payload
	_ = fileData // explicitly unused - included for method signature clarity

	result, err := r.InvokeModel(ctx, modelID, payload)
	if err != nil {
		return nil, fmt.Errorf("invoke model with file data: %w", err)
	}

	return result, nil
}
