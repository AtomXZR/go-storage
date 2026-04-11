package storage

import (
	"net/http"
	"path"
	"regexp"
	"strings"
)

//
//
//

var metadataRegEx *regexp.Regexp

func init() {
	metadataRegEx = regexp.MustCompile(`[^A-Za-z0-9\-_]+`)
}

//
//
//

// Normalize Metadate key to Canonical Header uwu;
func NormalizeMetadataKey(key string) string {
	key = metadataRegEx.ReplaceAllString(key, "")
	return http.CanonicalHeaderKey(key)
}

func NormalizeMetadata(metadata Metadata) Metadata {
	if metadata == nil {
		return nil
	}

	//

	out := make(Metadata, len(metadata))
	for k, v := range metadata {
		nk := NormalizeMetadataKey(k)
		out[nk] = v
	}

	//

	return out
}

//
//
//

func NormalizeKey(key string) (string, error) {
	nkey := strings.TrimSpace(key)

	if nkey == "" || nkey == "/" {
		return "", ErrInvalidKey.New(key)
	}

	//

	nkey = path.Clean(nkey)
	nkey = strings.TrimPrefix(nkey, "/")

	if nkey == "" || nkey == "." || nkey == ".." || strings.HasPrefix(nkey, "../") {
		return "", ErrInvalidKey.New(key)
	}

	return nkey, nil
}

//
//
//

func GetOptionsOrDefault(opts *GetOptions) *GetOptions {
	if opts == nil {
		return &GetOptions{}
	}
	return opts
}

func PutOptionsOrDefault(opts *PutOptions) *PutOptions {
	if opts == nil {
		return &PutOptions{
			ContentType: "application/octet-stream",
		}
	}

	newOpts := *opts

	//

	if newOpts.ContentType == "" {
		newOpts.ContentType = "application/octet-stream"
	}

	//

	return &newOpts
}

//
//
//
