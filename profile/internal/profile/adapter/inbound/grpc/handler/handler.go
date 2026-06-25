package handler

import (
	"context"

	profilev1 "github.com/anfastk/mergespace/contracts/gen/go/proto/profile/v1"

	"github.com/anfastk/mergespace/profile/internal/profile/application/dto"
	"github.com/anfastk/mergespace/profile/internal/profile/application/usecase"

	"connectrpc.com/connect"
)

type ProfileHandler struct {
	usecase *usecase.ProfileUseCase
}

func NewProfileHandler(
	usecase *usecase.ProfileUseCase,
) *ProfileHandler {

	return &ProfileHandler{
		usecase: usecase,
	}

}

func (h *ProfileHandler) UpdateProfile(ctx context.Context, req *connect.Request[profilev1.UpdateProfileRequest]) (*connect.Response[profilev1.UpdateProfileResponse], error) {

	request := &dto.UpdateProfileRequest{
		UserID: req.Msg.UserId,
	}

	if req.Msg.FirstName != nil {
		request.FirstName = req.Msg.FirstName
	}

	if req.Msg.LastName != nil {
		request.LastName = req.Msg.LastName
	}

	if req.Msg.Bio != nil {
		request.Bio = req.Msg.Bio
	}

	if req.Msg.AvatarUrl != nil {
		request.AvatarURL = req.Msg.AvatarUrl
	}

	if err := h.usecase.UpdateProfile(
		ctx,
		request,
	); err != nil {

		return nil, err
	}

	return connect.NewResponse(
		&profilev1.UpdateProfileResponse{
			Message: "profile updated successfully",
		},
	), nil

}
