package service

import (
	"encoding/json"
	"fmt"
)

// MemberManager: group contact member manager
type MemberManager struct {
	Group *User
}

// CreateMemberManagerFromGroupContact: create member manager by group contact info
func CreateMemberManagerFromGroupContact(session *Session, user *User) (*MemberManager, error) {
	b, err := WebWxBatchGetContact(session.WxWebCommon, session.WxWebXcg, session.Cookies, []*User{{
		EncryChatRoomId: user.EncryChatRoomId,
		UserName:        user.UserName,
	}})
	if err != nil {
		return nil, err
	}
	return CreateMemberManagerFromBytes(session, b)
}

// CreateMemberManagerFromBytes: create memeber manager by WxWebBatchGetContactResponse
func CreateMemberManagerFromBytes(session *Session, b []byte) (*MemberManager, error) {
	var gcr WxWebBatchGetContactResponse
	if err := json.Unmarshal(b, &gcr); err != nil {
		return nil, err
	}

	if gcr.BaseResponse.Ret != 0 {
		return nil, fmt.Errorf("WebWxBatchGetContact ret=%d", gcr.BaseResponse.Ret)
	}

	if gcr.ContactList == nil || gcr.Count < 1 || len(gcr.ContactList) < 1 {
		return nil, fmt.Errorf("ContactList empty")
	}

	mm := &MemberManager{
		Group: gcr.ContactList[0],
	}

	return mm, nil
}

// Update: get User details of group members
func (s *MemberManager) Update(session *Session) error {
	members := make([]*User, len(s.Group.MemberList))
	for i, v := range s.Group.MemberList {
		members[i] = &User{
			UserName:        v.UserName,
			EncryChatRoomId: s.Group.UserName,
		}
	}
	b, err := WebWxBatchGetContact(session.WxWebCommon, session.WxWebXcg, session.Cookies, members)
	if err != nil {
		return err
	}

	var gcr WxWebBatchGetContactResponse
	if err := json.Unmarshal(b, &gcr); err != nil {
		return err
	}
	s.Group.MemberList = gcr.ContactList
	return nil
}

// GetHeadImgUrlByGender: get head img url detail by gender
func (s *MemberManager) GetHeadImgUrlsByGender(sex int) []string {
	uris := make([]string, 0)
	for _, v := range s.Group.MemberList {
		if v.Sex == sex {
			uris = append(uris, v.HeadImgUrl)
		}
	}
	return uris
}

// GetContactsByGender: get contacts by gender
func (s *MemberManager) GetContactsByGender(sex int) []*User {
	contacts := make([]*User, 0)
	for _, v := range s.Group.MemberList {
		if v.Sex == sex {
			contacts = append(contacts, v)
		}
	}
	return contacts
}

// GetContactByUserName: get a certain member by username
func (s *MemberManager) GetContactByUserName(username string) *User {
	for _, v := range s.Group.MemberList {
		if v.UserName == username {
			return v
		}
	}
	return nil
}
