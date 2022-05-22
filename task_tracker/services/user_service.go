package services

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"sync"
)

type UserInfo struct {
	Id   uuid.UUID
	Name string
}

type UserInfoDataProvider interface {
	GetUserInfoChannel() <-chan UserInfo
}

type UserService interface {
	GetUser(userId uuid.UUID) (UserInfo, error)
	GetRandomUser() (UserInfo, error)
}

type randomizedSet struct {
	sync.Mutex
	usersIdx  map[uuid.UUID]int
	usersList []UserInfo
}

func NewUserService(userInfoDataProvider UserInfoDataProvider) UserService {
	rs := randomizedSet{usersIdx: make(map[uuid.UUID]int), usersList: make([]UserInfo, 0)}

	c := userInfoDataProvider.GetUserInfoChannel()
	go func() {
		for userInfoMsg := range c {
			userInfo := UserInfo{Id: userInfoMsg.Id, Name: userInfoMsg.Name}
			fmt.Println("Get new UserInfo: " + userInfoMsg.Id.String())
			rs.AddUser(userInfo)
		}
	}()

	return &rs
}

func (rs *randomizedSet) AddUser(user UserInfo) {
	rs.Lock()
	defer rs.Unlock()

	if idx, exist := rs.usersIdx[user.Id]; exist {
		rs.usersList[idx] = user
	} else {
		rs.usersList = append(rs.usersList, user)
		rs.usersIdx[user.Id] = len(rs.usersList) - 1
	}
}

func (rs *randomizedSet) GetUser(userId uuid.UUID) (UserInfo, error) {
	rs.Lock()
	defer rs.Unlock()

	if idx, exist := rs.usersIdx[userId]; exist {
		return rs.usersList[idx], nil
	} else {
		return UserInfo{}, fmt.Errorf("user with id=%s not found", userId.String())
	}
}

func (rs *randomizedSet) GetRandomUser() (UserInfo, error) {
	rs.Lock()
	defer rs.Unlock()

	_ = fmt.Errorf("not implemented")
	if len(rs.usersList) == 0 {
		return UserInfo{}, fmt.Errorf("there are no users can't give random")
	}

	idx := rand.Intn(len(rs.usersList))
	return rs.usersList[idx], nil
}

func (rs *randomizedSet) RemoveUser(userId uuid.UUID) bool {
	rs.Lock()
	defer rs.Unlock()

	if idx, exist := rs.usersIdx[userId]; exist {
		delete(rs.usersIdx, userId)

		// swap elements: element that we want to delete with last element
		// to have O(1) removal time complexity
		lastIdx := len(rs.usersList) - 1
		lastVal := rs.usersList[lastIdx]
		rs.usersList[idx] = lastVal
		rs.usersIdx[lastVal.Id] = idx
		//remove last element by truncating slice
		rs.usersList = rs.usersList[:lastIdx]
		return true
	} else {
		return false
	}
}
