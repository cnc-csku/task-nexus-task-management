package rest

import (
	"github.com/cnc-csku/task-nexus/task-management/domain/services"
)

type MemberHandler interface{}

type memberHandler struct {
	memberService services.MemberService
}

func NewMemberHandler(
	memberService services.MemberService,
) MemberHandler {
	return &memberHandler{
		memberService: memberService,
	}
}
