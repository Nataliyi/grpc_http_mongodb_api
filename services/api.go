package services

import (
	"api/consts"
	filter_models "api/models/filter"
	logger_models "api/models/logger"
	metric_models "api/models/metric"
	store_models "api/models/store"
	watcher_models "api/models/watcher"

	models "api/models/user"
	"api/util"
	"fmt"
	"time"

	"context"
	"net/http"

	pb "api/proto/gen/go"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Server for the gRPC API
type Server struct {
	// grpc server
	pb.UnimplementedUsersStoreServer
	// number of goroutines
	MaxProcessingGoroutines int
	// store client
	Store store_models.IStore
	// filter client
	Filter filter_models.IFilter
	// errors metric
	ErrorsMetric metric_models.IMetricCount
	// watcher channel
	WatcherCh watcher_models.WatcherChannel
	// logger
	Logger logger_models.ILogger
}

// Add new user to the store
func (s *Server) AddUser(ctx context.Context,
	request *pb.User) (resp *pb.UserResponse, err error) {
	// copy pb request to the user struct
	user := util.ConvertUserReq(request)
	// set updated and created time
	user.CreatedAt = models.CreatedAt(time.Now().UTC().Format(consts.TIME_FORMAT))
	user.UpdatedAt = models.UpdatedAt(time.Now().UTC().Format(consts.TIME_FORMAT))
	// using the loop for the duplicate key error
	for {
		// generate new uuid user id
		user.ID = util.GenID()
		// add new user to the store
		if err = s.Store.DoOne(store_models.ADD, user); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				// repeat insert
				continue
			}
			s.Logger.Error("AddUserError:", err.Error())
			// send to the errors metric
			s.ErrorsMetric.Add(1)
			// return error
			storeErr := fmt.Sprintf(consts.STORE_ERROR_FAILURE, err)
			return &pb.UserResponse{
				Status: http.StatusServiceUnavailable,
				Error:  &storeErr,
			}, nil
		}
		break
	}
	id := string(user.ID)
	s.Logger.Info("AddUser:", id)
	// inform
	s.WatcherCh <- fmt.Sprintf("ID: %s. Added new user.", id)
	// no errors
	return &pb.UserResponse{
		Id:     id,
		Status: http.StatusOK,
	}, nil
}

// Modify existed user by ID
func (s *Server) ModifyUser(ctx context.Context, request *pb.User) (*pb.UserResponse, error) {
	// check request
	if err := s.isValidRequest(request); err != nil {
		// return error
		storeErr := fmt.Sprintf(consts.STORE_BAD_REQUEST, err)
		return &pb.UserResponse{
			Id:     request.Id,
			Status: http.StatusBadRequest,
			Error:  &storeErr,
		}, nil
	}
	// copy pb request to the user struct
	user := util.ConvertUserReq(request)
	// set updated time
	user.UpdatedAt = models.UpdatedAt(time.Now().UTC().Format(consts.TIME_FORMAT))
	// modify the user in the store
	if err := s.Store.DoOne(store_models.MODIFY, user); err != nil {
		s.Logger.Error("ModifyUserError:", err.Error())
		// send to the errors metric
		s.ErrorsMetric.Add(1)
		// return error
		storeErr := fmt.Sprintf(consts.STORE_ERROR_FAILURE, err)
		return &pb.UserResponse{
			Id:     request.Id,
			Status: http.StatusServiceUnavailable,
			Error:  &storeErr,
		}, nil
	}
	s.Logger.Info("ModifyUser:", request.Id)
	// inform
	s.WatcherCh <- fmt.Sprintf("ID: %s. Modifyed user.", request.Id)
	// no errors
	return &pb.UserResponse{
		Id:     request.Id,
		Status: http.StatusOK,
	}, nil
}

// Delete an existed user by ID
func (s *Server) DeleteUser(ctx context.Context, request *pb.User) (*pb.UserResponse, error) {
	// check request
	if err := s.isValidRequest(request); err != nil {
		// return error
		storeErr := fmt.Sprintf(consts.STORE_BAD_REQUEST, err)
		return &pb.UserResponse{
			Id:     request.Id,
			Status: http.StatusBadRequest,
			Error:  &storeErr,
		}, nil
	}
	// copy pb request to the user struct
	user := util.ConvertUserReq(request)
	// delete the user in the store
	if err := s.Store.DoOne(store_models.DELETE, user); err != nil {
		s.Logger.Error("DeleteUserError:", err.Error())
		// send to the errors metric
		s.ErrorsMetric.Add(1)
		// return error
		storeErr := fmt.Sprintf(consts.STORE_ERROR_FAILURE, err)
		return &pb.UserResponse{
			Id:     request.Id,
			Status: http.StatusServiceUnavailable,
			Error:  &storeErr,
		}, nil
	}
	s.Logger.Info("DeleteUser:", request.Id)
	// inform
	s.WatcherCh <- fmt.Sprintf("ID: %s. Deleted user.", string(user.ID))
	// no errors
	return &pb.UserResponse{
		Id:     request.Id,
		Status: http.StatusOK,
	}, nil
}

// Get the list of all users
func (s *Server) GetAllUsers(
	context.Context, *emptypb.Empty) (*pb.UsersList, error) {
	// get all users
	results, respErr := s.Store.Get(store_models.GET_ALL, nil)
	// convert results to user
	users, err := util.ParseUsersToPb(results)
	// check parse error
	if err != nil {
		s.Logger.Error("GetAllUsersError:", err.Error())
		// send to the errors metric
		s.ErrorsMetric.Add(1)
		// return error
		storeErr := fmt.Sprintf(consts.STORE_ERROR_FAILURE, err)
		return &pb.UsersList{
			Status: http.StatusServiceUnavailable,
			Error:  &storeErr,
		}, nil
	}
	if respErr != nil {
		s.Logger.Error("GetAllUsersError:", err.Error())
		// send to the errors metric
		s.ErrorsMetric.Add(1)
		// return response with some errors
		storeErr := fmt.Sprintf(consts.STORE_ERROR_FAILURE, respErr)
		return &pb.UsersList{
			User:   users,
			Status: http.StatusAccepted,
			Error:  &storeErr,
		}, nil
	}
	// return response
	return &pb.UsersList{
		User:   users,
		Status: http.StatusOK,
	}, nil
}

// Get the list of filtered users
func (s *Server) GetUsers(
	ctx context.Context, filter *pb.UsersFilter) (*pb.UsersList, error) {
	// convert b request
	usersFilter := util.ConvertUserFilter(filter)
	// create new bson filter
	bsonUsersFilter := s.Filter.Filter(usersFilter)
	// get filtered users from the store
	results, respErr := s.Store.Get(store_models.GET_FILTERED, bsonUsersFilter)
	// convert results to user
	users, err := util.ParseUsersToPb(results)
	// check parse error
	if err != nil {
		s.Logger.Error("GetUsersError:", err.Error())
		// send to the errors metric
		s.ErrorsMetric.Add(1)
		// return error
		storeErr := fmt.Sprintf(consts.STORE_ERROR_FAILURE, err)
		return &pb.UsersList{
			Status: http.StatusServiceUnavailable,
			Error:  &storeErr,
		}, nil
	}
	if respErr != nil {
		s.Logger.Error("GetUsersError:", respErr.Error())
		// send to the errors metric
		s.ErrorsMetric.Add(1)
		// return response with some errors
		storeErr := fmt.Sprintf(consts.STORE_ERROR_FAILURE, respErr)
		return &pb.UsersList{
			User:   users,
			Status: http.StatusAccepted,
			Error:  &storeErr,
		}, nil
	}
	// return response
	return &pb.UsersList{
		User:   users,
		Status: http.StatusOK,
	}, nil
}

// Check the valid request
func (s *Server) isValidRequest(request *pb.User) error {
	if request.Id == "" {
		return fmt.Errorf("the id field not set")
	}
	if request.CreatedAt != "" {
		return fmt.Errorf("created_at must not set")
	}
	if request.UpdatedAt != "" {
		return fmt.Errorf("updated_at must not set")
	}
	return nil
}
