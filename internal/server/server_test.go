package server_test

import (
	"context"
	"testing"
	"testproject/internal/server"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	_, err := server.NewServer()
	assert.Nil(t, err)
}

func TestServer_StartNoTls(t *testing.T) {
	s, err := server.NewServer()
	assert.Nil(t, err)
	closing := false
	go func() {
		err := s.Start(false)
		if !closing {
			assert.Nil(t, err)
		}
	}()
	time.Sleep(time.Second * 5)
	closing = true
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	assert.Nil(t, s.Stop(ctx))
}

func TestServer_StartTls(t *testing.T) {
	s, err := server.NewServer()
	assert.Nil(t, err)
	closing := false
	go func() {
		err := s.Start(true)
		if !closing {
			assert.Nil(t, err)
		}
	}()
	time.Sleep(time.Second * 5)
	closing = true
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	assert.Nil(t, s.Stop(ctx))
}
