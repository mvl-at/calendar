package calendar

import (
	"encoding/json"
	"fmt"
	"github.com/mvl-at/model"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

type UserInfo struct {
	Member *model.Member `json:"member"`
	Roles  []*model.Role `json:"roles"`
}

func fetchEvents(events *[]*model.Event, from string, to string) {
	var jsonData []byte
	resp, err := http.Get(fmt.Sprintf("%s/eventsrange?from=%s&to=%s", conf.RestHost, from, to))

	if err != nil {
		errLogger.Println(err.Error())
		return
	}
	if resp.Body != nil {
		jsonData, _ = ioutil.ReadAll(resp.Body)
	} else {
		return
	}
	json.Unmarshal(jsonData, events)
	eventsVal := *events
	sort.Slice(eventsVal, func(i, j int) bool {
		iD := eventsVal[i].Date
		iT := eventsVal[i].Time
		iTime := time.Date(iD.Year(), iD.Month(), iD.Day(), iT.Hour(), iT.Minute(), 0, 0, time.Local)
		jD := eventsVal[j].Date
		jT := eventsVal[j].Time
		jTime := time.Date(jD.Year(), jD.Month(), jD.Day(), jT.Hour(), jT.Minute(), 0, 0, time.Local)
		return iTime.Unix() < jTime.Unix()
	})
}

func fetchObmAndKpm() (obm model.Member, kpm model.Member) {
	obm = model.Member{}
	kpm = model.Member{}
	var jsonData []byte
	resp, err := http.Get(conf.RestHost + "/leaderRolesMembers")
	if err != nil {
		errLogger.Println(err.Error())
		return
	}
	if resp.Body != nil {
		jsonData, _ = ioutil.ReadAll(resp.Body)
	} else {
		return
	}
	leaders := make([]model.LeaderRoleMember, 0)
	json.Unmarshal(jsonData, &leaders)
	for _, leader := range leaders {
		if strings.ToLower(leader.LeaderRole.Name) == strings.ToLower(conf.Obm.Role) && !leader.Deputy {
			obm = *leader.Member
		}
		if strings.ToLower(leader.LeaderRole.Name) == strings.ToLower(conf.Kpm.Role) && !leader.Deputy {
			kpm = *leader.Member
		}
	}
	return
}

func fetchAuthor(jwt string) (author UserInfo, ok bool) {
	ok = false
	author = UserInfo{}
	request, err := http.NewRequest(http.MethodGet, conf.RestHost+"/userinfo", nil)
	if err != nil {
		errLogger.Println(err.Error())
	} else {
		client := http.DefaultClient
		request.Header.Set("Access-token", jwt)
		response, err := client.Do(request)
		if err != nil {
			errLogger.Println(err.Error())
		} else {
			decoder := json.NewDecoder(response.Body)
			decoder.Decode(&author)
			ok = author.Member != nil
		}
	}
	return
}

func hasRole(user *UserInfo, role string) bool {
	r := strings.ToLower(role)
	if r == "root" {
		return true
	}
	for _, ro := range user.Roles {
		if strings.ToLower(ro.Id) == r {
			return true
		}
	}
	return false
}
