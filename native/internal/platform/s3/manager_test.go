package s3

import (
	"context"
	"errors"
	"testing"

	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type fakeBucketClient struct {
	headErrors  []error
	headCalls   int
	createError error
	createCalls int
}

func (f *fakeBucketClient) HeadBucket(context.Context, *awss3.HeadBucketInput, ...func(*awss3.Options)) (*awss3.HeadBucketOutput, error) {
	index := f.headCalls
	f.headCalls++
	if index < len(f.headErrors) && f.headErrors[index] != nil {
		return nil, f.headErrors[index]
	}
	return &awss3.HeadBucketOutput{}, nil
}

func (f *fakeBucketClient) CreateBucket(context.Context, *awss3.CreateBucketInput, ...func(*awss3.Options)) (*awss3.CreateBucketOutput, error) {
	f.createCalls++
	if f.createError != nil {
		return nil, f.createError
	}
	return &awss3.CreateBucketOutput{}, nil
}

func TestEnsureBucket(t *testing.T) {
	tests := []struct {
		name        string
		client      *fakeBucketClient
		wantErr     bool
		createCalls int
	}{
		{name: "already exists", client: &fakeBucketClient{}, createCalls: 0},
		{name: "creates missing bucket", client: &fakeBucketClient{headErrors: []error{errors.New("missing")}}, createCalls: 1},
		{name: "accepts concurrent creation", client: &fakeBucketClient{headErrors: []error{errors.New("missing"), nil}, createError: errors.New("already created")}, createCalls: 1},
		{name: "returns create and recheck failure", client: &fakeBucketClient{headErrors: []error{errors.New("missing"), errors.New("still missing")}, createError: errors.New("create failed")}, wantErr: true, createCalls: 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ensureBucket(context.Background(), test.client, "quick-admin")
			if (err != nil) != test.wantErr {
				t.Fatalf("ensureBucket() error = %v, wantErr %v", err, test.wantErr)
			}
			if test.client.createCalls != test.createCalls {
				t.Fatalf("create calls = %d, want %d", test.client.createCalls, test.createCalls)
			}
		})
	}
}
