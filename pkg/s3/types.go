package s3

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	transport "github.com/aws/smithy-go/endpoints"
)

var _ aws.CredentialsProvider = credentialProvider{}
var _ s3.EndpointResolverV2 = endpointResolver{}

type S3Client struct {
	// This is the S3 client.
	bucket        string
	client        *s3.Client
	encryptionKey string
}

type S3ClientOpts struct {
	// This is the S3 client options.
	Bucket        string
	Region        string
	AccessKey     string
	SecretKey     string
	EndpointURL   string
	EncryptionKey string
}

func NewS3Client(ctx context.Context, opts *S3ClientOpts) (*S3Client, error) {
	// NewS3Client creates a new S3 client
	var awsOptFns []func(*config.LoadOptions) error
	var s3OptFns []func(*s3.Options)

	if opts.AccessKey != "" && opts.SecretKey != "" {
		awsOptFns = append(awsOptFns, config.WithCredentialsProvider(credentialProvider{accesKey: opts.AccessKey, secretKey: opts.SecretKey}))
	}

	if opts.EndpointURL != "" {
		endpoint := endpointResolver{
			endpoint: opts.EndpointURL,
		}

		s3OptFns = append(s3OptFns, s3.WithEndpointResolverV2(endpoint))
	}

	if opts.Region != "" {
		awsOptFns = append(awsOptFns, config.WithRegion(opts.Region))
	}

	cfg, err := config.LoadDefaultConfig(ctx, awsOptFns...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	if opts.Bucket == "" {
		return nil, fmt.Errorf("bucket is required")
	}

	s3Client := s3.NewFromConfig(cfg, s3OptFns...)

	return &S3Client{
		bucket:        opts.Bucket,
		client:        s3Client,
		encryptionKey: opts.EncryptionKey,
	}, nil

}

type credentialProvider struct {
	accesKey  string
	secretKey string
}

func (p credentialProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     p.accesKey,
		SecretAccessKey: p.secretKey,
	}, nil
}

type endpointResolver struct {
	endpoint string
}

func (r endpointResolver) ResolveEndpoint(_ context.Context, params s3.EndpointParameters) (transport.Endpoint, error) {
	u, err := url.Parse(r.endpoint)
	if err != nil {
		return transport.Endpoint{}, fmt.Errorf("failed to parse endpoint URL: %w", err)
	}

	u.Path += "/" + *params.Bucket

	return transport.Endpoint{
		URI: *u,
	}, nil
}
