/*package storage package to store the data in mongo, this package
creates a session and persists the session, to contact with mongo.
And creates document in mongo.
*/
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

//getSession to get mongo session, this function validates the session is empty or not
//if session is valid then it returns the mongo session, otherwise creates a mongo session
//and returns the copy session.
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

//ChatItem to store the coonversation data between user and AI
type ChatItem struct {
	UserInput string `bson:"userInput"`
	AIOutPut  string `bson:"AIOutput"`
}

//formJSON takes input as string array and forms the chatItem array
func formJSON(chathistory []string) []ChatItem {
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

//AddToMongo is the interface to communicate to mongo by mongo session
//it creates document in mongo
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
