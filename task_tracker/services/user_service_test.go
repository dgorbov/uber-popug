package services

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type userInfoProviderStub struct {
	userInfoChannel chan UserInfo
}

func (stub *userInfoProviderStub) GetUserInfoChannel() <-chan UserInfo {
	return stub.userInfoChannel
}

func createUserServiceWithStubProvider(bufferSize int) (UserService, userInfoProviderStub) {
	provider := userInfoProviderStub{userInfoChannel: make(chan UserInfo, bufferSize)}
	return NewUserService(&provider), provider
}

func TestUserService_GetUser_InjectProviderAndProvideValidUserId_UserReceived(t *testing.T) {
	t.Skip("requires some work...")
	userService, provider := createUserServiceWithStubProvider(0)
	user := UserInfo{Id: uuid.New(), Name: "test user"}

	provider.userInfoChannel <- user
	addedUser, _ := userService.GetUser(user.Id)

	assert.EqualValues(t, user, addedUser)
}

func TestRandomizedSet_GetRandomUser_ProvideFewUsers_UserReceived(t *testing.T) {
	t.Skip("requires some work...")

	userService, provider := createUserServiceWithStubProvider(3)
	provider.userInfoChannel <- UserInfo{Id: uuid.New(), Name: "test user 1"}
	provider.userInfoChannel <- UserInfo{Id: uuid.New(), Name: "test user 2"}
	provider.userInfoChannel <- UserInfo{Id: uuid.New(), Name: "test user 3"}

	randomUser, err := userService.GetRandomUser()
	fmt.Println("got random user with name: " + randomUser.Name)
	assert.Nil(t, err)
}

func TestRandomizedSet_RemoveUser_AddRemoveUsers_CanGetOnlyExisting(t *testing.T) {
	userService, _ := createUserServiceWithStubProvider(3)
	rs := userService.(*randomizedSet)

	user1 := UserInfo{Id: uuid.New(), Name: "test user 1"}
	user2 := UserInfo{Id: uuid.New(), Name: "test user 2"}

	rs.AddUser(user1)
	rs.AddUser(user2)
	rs.RemoveUser(user1.Id)
	_, err := rs.GetUser(user1.Id)
	assert.NotNil(t, err)
	_, err = rs.GetUser(user2.Id)
	assert.Nil(t, err)

	user3 := UserInfo{Id: uuid.New(), Name: "test user 3"}
	rs.AddUser(user3)
	_, err = rs.GetUser(user1.Id)
	assert.NotNil(t, err)
	_, err = rs.GetUser(user2.Id)
	assert.Nil(t, err)
	_, err = rs.GetUser(user3.Id)
	assert.Nil(t, err)
}
