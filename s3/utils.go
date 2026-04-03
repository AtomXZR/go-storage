package s3

import (
	"net/http"
	"strings"

	"github.com/AtomXZR/go-storage"
	"github.com/minio/minio-go/v7"
)

func parseMetadata(headers http.Header) storage.Metadata {
	result := make(storage.Metadata)

	for k, v := range headers {
		nk := storage.NormalizeMetadataKey(k)

		nk = strings.TrimPrefix(nk, "X-Amz-Meta-")

		if len(v) == 0 {
			continue
		}

		result[nk] = v[0]
	}

	return result
}

func objectInfoToStats(objectInfo minio.ObjectInfo) storage.Stats {
	return storage.Stats{
		ETag:         objectInfo.ETag,
		ContentType:  objectInfo.ContentType,
		Size:         objectInfo.Size,
		LastModified: objectInfo.LastModified,

		Metadata: parseMetadata(objectInfo.Metadata),
	}
}
