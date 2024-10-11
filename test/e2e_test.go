package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/efficientgo/core/testutil"
	"github.com/efficientgo/e2e"
	e2edb "github.com/efficientgo/e2e/db"
	e2einteractive "github.com/efficientgo/e2e/interactive"
	e2emonitoring "github.com/efficientgo/e2e/monitoring"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/thanos-io/objstore/providers/s3"
)

const (
	MinioAccessKey = "Cheescake"
	MinioSecretKey = "supersecret"
)

func uploadTestInput(client *minio.Client, objectName string, size int64) error {
	data := bytes.Repeat([]byte("x"), int(size)) // Create a virtual file of size 'size'

	// Upload file to MinIO
	_, err := client.PutObject(context.Background(), "test", objectName, bytes.NewReader(data), size, minio.PutObjectOptions{
		ContentType: "application/text", // Set the appropriate content type
	})
	return err
}

func TestLabeler_LabelObject(t *testing.T) {
	t.Parallel()
	// Create Docker environment
	e, err := e2e.NewDockerEnvironment("labeler")
	testutil.Ok(t, err)
	t.Cleanup(e.Close)

	// Start monitoring
	mon, err := e2emonitoring.Start(e)
	testutil.Ok(t, err)
	testutil.Ok(t, mon.OpenUserInterfaceInBrowser())

	// Create MinIO container
	minioContainer := e2edb.NewMinio(e, "mintest", "bkt")
	testutil.Ok(t, e2e.StartAndWaitReady(minioContainer))

	// Create MinIO client
	minioClient, err := minio.New(minioContainer.InternalEndpoint(e2edb.AccessPortName), &minio.Options{
		Creds:  credentials.NewStaticV4(MinioAccessKey, MinioSecretKey, ""),
		Secure: false,
	})
	testutil.Ok(t, err)

	// Create bucket if it does not exist
	err = minioClient.MakeBucket(context.Background(), "test", minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), "test")
		if errBucketExists == nil && exists {
			log.Printf("Bucket already exists: %s", "test")
		} else {
			log.Fatalf("Failed to create bucket: %v", err)
		}
	}

	// Configure S3
	config := struct {
		Type   string
		Config s3.Config
	}{
		Type: "S3",
		Config: s3.Config{
			Bucket:    "test",
			AccessKey: MinioAccessKey,
			SecretKey: MinioSecretKey,
			Endpoint:  minioContainer.InternalEndpoint(e2edb.AccessPortName),
			Insecure:  true,
		},
	}

	// Serialize config
	configBytes, err := json.Marshal(config)
	if err != nil {
		log.Fatalf("failed to marshal config: %v", err)
	}

	// Upload test input
	testutil.Ok(t, uploadTestInput(minioClient, "object.txt", 2e6))

	// Start Labeler
	labeler := e.Runnable("labeler").
		WithPorts(map[string]int{"http": 8080}).
		Init(e2e.StartOptions{
			Image:     "labeler:test",
			LimitCPUs: 4.0,
			Command: e2e.NewCommand(
				"labeler",
				"-listen-address=:8080",
				"-objstore.config="+string(configBytes),
			),
		})

	testutil.Ok(t, e2e.StartAndWaitReady(labeler))

	// Build URL (ensure object_id matches uploaded object's name)
	url := fmt.Sprintf("http://%s/label_object?object_id=object.txt", labeler.InternalEndpoint("http"))

	// k6 test script
	k6Command := fmt.Sprintf(` 
		import http from 'k6/http';
		import { check, sleep } from 'k6';

		export default function() {
			const res = http.get('%s');
			check(res, {
				'is status 200': (r) => r.status === 200,
				'response': (r) => r.body.includes('{"object_id": "object.txt", "sum":6221600000, "checksum":"SUUr1234567890abcdef"}'),
			});
			sleep(0.5);
		}
	`, url)

	// Start k6
	k6 := e.Runnable("k6").Init(e2e.StartOptions{
		Command: e2e.NewCommandRunUntilStop(),
		Image:   "grafana/k6:0.39.0",
	})
	testutil.Ok(t, e2e.StartAndWaitReady(k6))

	// Execute k6 test
	testutil.Ok(t, k6.Exec(e2e.NewCommand("/bin/sh", "-c", k6Command)))
	testutil.Ok(t, e2einteractive.RunUntilEndpointHit())
}
