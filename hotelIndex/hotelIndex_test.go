package hotelindex

import (
	"goapp/role"
	"log/slog"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestHotelUriGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockHotelIndexInterface(ctrl)

	mockRoll := role.NewMockRoleInterface(ctrl)
	mockRoll.EXPECT().GetUserRoleByFacility("any", "any").
		Return(&role.Roles{
			HotelID: "hoteltestID",
			Role:    "staff",
		}, nil).
		Times(1)
	mockRepo.EXPECT().
		GetHotelIndexInfo("testString", "TestString2").
		Return(&HotelIndex{
			HotelName: "test",
			Authority: "staff",
		}, nil).
		Times(1)
	uc := NewHotelIndexHandler(mockRepo, mockRoll)
	_, err := uc.RoleInterfaces.GetUserRoleByFacility("any", "any")
	if err != nil {
		slog.Error("テストに失敗しました", "error", err)
		return
	}
	_, err = uc.HotelInterface.GetHotelIndexInfo("testString", "TestString2")
	if err != nil {
		slog.Error("テストに失敗しました。", "error", err)
		return
	}
	slog.Info("テストが完了しました！")
}
