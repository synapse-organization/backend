package modules

import (
	"barista/pkg/errors"
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/repo"
	"barista/pkg/utils"
	"context"
	"math/rand"
	"strconv"
	"time"
)

const (
	commentsLimit = 5
)

type CafeHandler struct {
	CafeRepo    repo.CafesRepo
	Rating      repo.RatingsRepo
	CommentRepo repo.CommentsRepo
	ImageRepo   repo.ImageRepo
	EventRepo   repo.EventRepo
	UserRepo	repo.UsersRepo
}

func (c CafeHandler) Create(ctx context.Context, cafe *models.Cafe) error {
	err := c.CafeRepo.Create(ctx, cafe)
	if err != nil {
		log.GetLog().Errorf("Unable to create cafe. error: %v", err)
		return err
	}
	return nil
}

func (c CafeHandler) GetCafes() {
	panic("implement me")
}

func (c CafeHandler) GetCafeByID() {
	panic("implement me")
}

func (c CafeHandler) SearchCafe(ctx context.Context, name string, address string, location string, category string) ([]models.Cafe, error) {
	cafes, err := c.CafeRepo.SearchCafe(ctx, name, address, location, category)
	if err != nil {
		log.GetLog().Errorf("Unable to search cafe. error: %v", err)
		return nil, err
	}
	for i, cafe := range cafes {
		cafes[i].Rating, err = c.Rating.GetCafesRating(ctx, cafe.ID)
		if err != nil {
			log.GetLog().Errorf("Unable to get cafe rating. error: %v", err)
		}
	}

	for i := range cafes {
		images, err := c.ImageRepo.GetByCafeID(ctx, cafes[i].ID)
		if err != nil {
			log.GetLog().Errorf("Unable to get cafe images. error: %v", err)
			continue
		}

		comments, err := c.CommentRepo.GetLast(ctx, cafes[i].ID, commentsLimit, 0)
		if err != nil {
			log.GetLog().Errorf("Unable to get cafe comments. error: %v", err)
			continue
		}

		for _, image := range images {
			cafes[i].Images = append(cafes[i].Images, image.ID)
		}

		for _, comment := range comments {
			cafes[i].Comments = append(cafes[i].Comments, *comment)
		}
	}

	return cafes, err
}

// func (c CafeHandler) PublicCafeProfile(ctx context.Context, cafeID string) (*models.Cafe, error) {
// 	cafe_id, err := strconv.Atoi(cafeID)
// 	if err != nil {
// 		log.GetLog().Errorf("Unable to convert userID to int32. error: %v", err)
// 		return nil, err
// 	}

// 	cafe, err := c.CafeRepo.GetByID(ctx, int32(cafe_id))
// 	if err != nil {
// 		log.GetLog().Errorf("Cafe id does not exist. error: %v", err)
// 		return nil, err
// 	}

// 	return cafe, nil
// }

func (c CafeHandler) AddComment(ctx context.Context, cafeID int32, userID string, comment string) (CommentWithUserName, error) {
	CWU := CommentWithUserName{}

	user_id, err := strconv.Atoi(userID)
	if err != nil {
		log.GetLog().Errorf("Unable to convert user id to int32. error: %v", err)
		return CWU, err
	}

	commentID := rand.Int31()

	AddedComment := &models.Comment{
		ID:      commentID,
		UserID:  int32(user_id),
		CafeID:  cafeID,
		Comment: comment,
		Date:    time.Now().UTC(),
	}

	err = c.CommentRepo.Create(ctx, AddedComment)
	if err != nil {
		log.GetLog().Errorf("Unable to add comment. error: %v", err)
		return CWU, err
	}

	user, err := c.UserRepo.GetByID(ctx, int32(user_id))
	if err != nil {
		log.GetLog().Errorf("Unable to get user name for user id %d: %v", user_id, err)
		return CWU, err
	}

	CWU.Comment = AddedComment
	CWU.UserFirstName = user.FirstName
	CWU.UserLastName = user.LastName

	return CWU, err
}

type CommentWithUserName struct {
	*models.Comment
	UserFirstName string `json:"user_first_name"`
	UserLastName	string `json:"user_last_name"`
}

func (c CafeHandler) GetComments(ctx context.Context, cafeID int32, counter int) ([]CommentWithUserName, error) {
	offset := (counter - 1) * commentsLimit
	comments, err := c.CommentRepo.GetLast(ctx, cafeID, commentsLimit, offset)
	if err != nil {
		log.GetLog().Errorf("Unable to get 5 last comments. error: %v", comments)
		return nil, err
	}

    var commentsWithNames []CommentWithUserName

    for _, comment := range comments {
        userName, err := c.UserRepo.GetByID(ctx, comment.UserID)
        if err != nil {
            log.GetLog().Errorf("Unable to get user name for user ID %d: %v", comment.UserID, err)
            return nil, err
        }

        commentWithUserName := CommentWithUserName{
            Comment:  comment,
            UserFirstName: userName.FirstName,
			UserLastName: userName.LastName,
        }

        commentsWithNames = append(commentsWithNames, commentWithUserName)
    }

	return commentsWithNames, err
}

func (c CafeHandler) CreateEvent(ctx context.Context, cafeID int32, name string, description string, startTime time.Time, endTime time.Time, imageID string) error {
	start_time := startTime.UTC()
	end_time := endTime.UTC()

	if !utils.CheckStartTime(start_time) {
		return errors.ErrStartTimeInvalid.Error()
	}

	if !utils.CheckEndTime(start_time, end_time) {
		return errors.ErrEndTimeInvalid.Error()
	}

	exists, err := c.ImageRepo.CheckExistence(ctx, imageID)
	if err != nil {
		log.GetLog().Errorf("Unable to check image existence. error: %v", err)
		return err
	}

	if !exists {
		log.GetLog().Errorf("image doesn't exist")
		return errors.ErrImageInvalid.Error()
	}

	eventID := rand.Int31()
	newEvent := &models.Event{
		ID:          eventID,
		CafeID:      cafeID,
		Name:        name,
		Description: description,
		StartTime:   start_time,
		EndTime:     end_time,
		ImageID:     imageID,
	}
	
	err = c.EventRepo.CreateEventForCafe(ctx, newEvent)
	if err != nil {
		log.GetLog().Errorf("Unable to create event for cafe. error: %v", err)
		return err
	}

	return err
}
