package util

import (
	"encoding/json"
	"regexp"

	user_models "api/models/user"
	pb "api/proto/gen/go"

	store_models "api/models/store"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// matched origin
func AllowedOrigin(origin string) bool {
	if viper.GetString("cors") == "*" {
		return true
	}
	if matched, _ := regexp.MatchString(viper.GetString("cors"), origin); matched {
		return true
	}
	return false
}

/*
Convert pb req to the user.
*/
func ConvertUserReq(pbReq *pb.User) *user_models.User {

	// convert the req to the map
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(pbReq)
	json.Unmarshal(inrec, &inInterface)

	// set id key
	if inInterface["id"] != nil {
		inInterface["_id"] = inInterface["id"]
		delete(inInterface, "id")
	}

	// convert the map to the user
	user := user_models.User{}
	jsonbody, _ := json.Marshal(inInterface)
	json.Unmarshal(jsonbody, &user)

	return &user

}

// convert pb UserFilter
func ConvertUserFilter(pbReq *pb.UsersFilter) (inInterface map[string][]interface{}) {

	// convert the request to the map
	inrec, _ := json.Marshal(pbReq)
	json.Unmarshal(inrec, &inInterface)

	// set id key
	if inInterface["id"] != nil {
		inInterface["_id"] = inInterface["id"]
		delete(inInterface, "id")
	}

	return
}

// Generate new users uuid ID
func GenID() user_models.ID {
	uuidID := uuid.Must(uuid.NewRandom()).String()
	return user_models.ID(uuidID)
}

// Parse users from the get-response to the users list
func ParseUsersTo(
	results []store_models.IStoreGetResponse) ([]*user_models.User, error) {

	usersResults := make([]*user_models.User, 0)

	// decode results
	for _, res := range results {
		jsonString, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}
		u := user_models.User{}
		if err := json.Unmarshal(jsonString, &u); err != nil {
			return nil, err
		}
		usersResults = append(usersResults, &u)
	}
	return usersResults, nil
}

// Parse users from the get-response to the pb users list
func ParseUsersToPb(
	results []store_models.IStoreGetResponse) ([]*pb.User, error) {

	usersResults := make([]*pb.User, 0)

	// decode results
	for _, res := range results {
		// set id key
		if res["_id"] != nil {
			res["id"] = res["_id"]
			delete(res, "_id")
		}
		jsonString, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}
		u := pb.User{}
		if err := json.Unmarshal(jsonString, &u); err != nil {
			return nil, err
		}
		usersResults = append(usersResults, &u)
	}
	return usersResults, nil
}
