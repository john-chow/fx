package core

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/metrue/fx/pkg/docker/mocks"
)

func TestIsFxDockerImageReady(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockCli := mock_docker.NewMockIClient(mockCtrl)
	mockCli.EXPECT().ListImagesWithReference(gomock.Any()).Return([]string{"mock image repo"}, nil).AnyTimes()

	ready := IsFxDockerImageReady(mockCli)
	if ready == false {
		t.Fatalf("should not be ready")
	}
}

func TestPullBaseImage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockCli := mock_docker.NewMockIClient(mockCtrl)
	mockCli.EXPECT().Pull(gomock.Any()).Return(nil).AnyTimes()

	PullBaseIamage(mockCli)
}
