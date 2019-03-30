package storage

import (
	"fmt"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

var session *mgo.Session

type chatHistory struct {
	SessionID   time.Time  `bson:"sessionID"`
	ChatHistory []ChatItem `bson:"chat_history"`
}

var c sync.Mutex

func getSession() (*mgo.Session, error) {
	var err error
	if session != nil {
		return session.Copy(), err
	}
	c.Lock()
	if session == nil {
		session, err = mgo.Dial("localhost")
		if err != nil {
			return session, err
		}
	}
	c.Unlock()
	return session.Copy(), err
}

//[{userINput:"", AIOutput:""}, ....]
type ChatItem struct {
	UserInput string `bson:"userInput"`
	AIOutPut  string `bson:"AIOutput"`
}

func formJson(chathistory []string) []ChatItem {
	var chatHistoryData []ChatItem
	var chatData ChatItem
	for index, data := range chathistory {
		if index%2 == 1 || index == 1 {
			chatData.AIOutPut = data
		} else if index == 0 || index%2 == 0 {
			chatData.UserInput = data
		}
		if (index == 1 || index%2 == 1) && index != 0 {
			chatHistoryData = append(chatHistoryData, chatData)
			chatData = ChatItem{}
		}
	}
	return chatHistoryData
}

func AddToMongo(sessionID time.Time, chatArray []string) error {
	session, err := getSession()
	if err != nil {
		fmt.Println("error getting session", err)
		return err
	}
	defer session.Close()
	c := session.DB("texa").C("chat")
	return c.Insert(chatHistory{sessionID, formJson(chatArray)})
}
