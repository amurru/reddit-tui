package images

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"path"
	"reddittui/client/common"
	"strings"
	"time"
)

const defaultMaxImageSizeBytes = 20 << 20 // 20 MB

type PreviewClient struct {
	httpClient *http.Client
	renderer   Renderer
	maxBytes   int64
	protocol   RenderProtocol
}

func NewPreviewClient(timeoutSeconds int, protocol string) PreviewClient {
	return PreviewClient{
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
		renderer: NewTermimgRenderer(),
		maxBytes: defaultMaxImageSizeBytes,
		protocol: RenderProtocol(protocol),
	}
}

func (p PreviewClient) RenderFromURL(imageURL string, width, height int) (string, error) {
	parsed, err := url.Parse(imageURL)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrUnsupportedImageURL, imageURL)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("%w: %s", ErrUnsupportedImageURL, imageURL)
	}

	req, err := http.NewRequest(http.MethodGet, imageURL, nil)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrCannotFetchImage, err)
	}
	req.Header.Add(common.UserAgentHeaderKey, common.UserAgentHeaderValue)

	res, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrCannotFetchImage, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: status code %d", ErrCannotFetchImage, res.StatusCode)
	}

	contentType := strings.ToLower(strings.TrimSpace(strings.Split(res.Header.Get("Content-Type"), ";")[0]))
	if !strings.HasPrefix(contentType, "image/") && !looksLikeImagePath(parsed.Path) {
		return "", fmt.Errorf("%w: %s", ErrUnsupportedImageURL, imageURL)
	}

	body, err := io.ReadAll(io.LimitReader(res.Body, p.maxBytes+1))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrCannotFetchImage, err)
	}

	if int64(len(body)) > p.maxBytes {
		return "", ErrImageTooLarge
	}

	img, _, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrCannotDecodeImage, err)
	}

	rendered, err := p.renderer.Render(img, RenderOptions{
		Width:    width,
		Height:   height,
		Protocol: p.protocol,
	})
	if err != nil {
		return "", err
	}

	return rendered, nil
}

func looksLikeImagePath(imagePath string) bool {
	ext := strings.ToLower(path.Ext(imagePath))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".tiff", ".webp":
		return true
	default:
		return false
	}
}
