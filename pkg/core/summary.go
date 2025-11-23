package core

import (
	"context"
	"fmt"
	"unicode/utf8"
)

func (s *Service) Summary(ctx context.Context, url string) (*Response, error) {
	page, err := s.crawler.FetchPage(ctx, url)
	if err != nil {
		fmt.Errorf("failed to fetch page: %w", err)
	}

	if utf8.RuneCountInString(page) <= 100 {
		return &Response{Message: page}, nil
	}

	runes := []rune(page)
	return &Response{Message: string(runes[:100])}, nil
}
