package mongo

import (
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MemberRepositoryImpl struct {
	mongoClient *mongo.Client
}

func NewMemberRepository(mongoClient *mongo.Client) repositories.MemberRepository {
	return &MemberRepositoryImpl{mongoClient}
}

func (r *MemberRepositoryImpl) GetMembers() {}
