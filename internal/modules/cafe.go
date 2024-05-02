package modules

import (
	"barista/pkg/errors"
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/repo"
	"barista/pkg/utils"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	commentsLimit = 5
)

type CafeHandler struct {
	CafeRepo        repo.CafesRepo
	Rating          repo.RatingsRepo
	CommentRepo     repo.CommentsRepo
	ImageRepo       repo.ImageRepo
	EventRepo       repo.EventRepo
	UserRepo        repo.UsersRepo
	ReservationRepo repo.ReservationRepo
	MenuItemRepo    repo.MenuItemsRepo
}

func (c CafeHandler) Create(ctx context.Context, cafe *models.Cafe) error {
	cafeID, err := c.CafeRepo.Create(ctx, cafe)
	if err != nil {
		log.GetLog().Errorf("Unable to create cafe. error: %v", err)
		return err
	}
	for _, photoID := range cafe.Images {
		err = c.ImageRepo.Create(ctx, &models.Image{
			ID:     photoID,
			CafeID: cafeID,
		})
	}

	return err
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

type PublicCafeProvinceCity struct {
	ID           int32                 `json:"id"`
	OwnerID      int32                 `json:"owner_id"`
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	OpeningTime  int8                  `json:"opening_time"`
	ClosingTime  int8                  `json:"closing_time"`
	Comments     []models.Comment      `json:"comments"`
	Rating       float64               `json:"rating"`
	Images       []string              `json:"photos"`
	Events       []models.Event        `json:"events"`
	Capacity     int32                 `json:"capacity"`
	ContactInfo  models.ContactInfo    `json:"contact_info"`
	Categories   []models.CafeCategory `json:"categories"`
	ProvinceName string                `json:"province_name"`
	CityName     string                `json:"city_name"`
}

func (c CafeHandler) PublicCafeProfile(ctx context.Context, cafeID int32) (*PublicCafeProvinceCity, error) {
	cafe, err := c.CafeRepo.GetByID(ctx, int32(cafeID))
	if err != nil {
		log.GetLog().Errorf("Cafe id does not exist. error: %v", err)
		return nil, err
	}

	comments, err := c.CommentRepo.GetLast(ctx, int32(cafeID), commentsLimit, 0)
	if err != nil {
		log.GetLog().Errorf("Unable to get last %v comments. error: %v", commentsLimit, err)
		return nil, err
	}

	cafe.Comments = make([]models.Comment, len(comments))
	for i, comment := range comments {
		cafe.Comments[i] = *comment
	}

	events, err := c.EventRepo.GetEventsByCafeID(ctx, int32(cafeID))
	if err != nil {
		log.GetLog().Errorf("Unable to get events by cafe id. error: %v", err)
		return nil, err
	}

	cafe.Events = make([]models.Event, len(events))
	for i, event := range events {
		cafe.Events[i] = *event
	}

	cafe.Rating, err = c.Rating.GetCafesRating(ctx, int32(cafeID))
	if err != nil {
		log.GetLog().Errorf("Unable to get rating by cafe id. error: %v", err)
		return nil, err
	}

	photos, err := c.ImageRepo.GetByCafeID(ctx, int32(cafeID))
	if err != nil {
		log.GetLog().Errorf("Unable to get photos by cafe id. error: %v", err)
		return nil, err
	}

	cafe.Images = make([]string, len(photos))
	for i, photo := range photos {
		cafe.Images[i] = photo.ID
	}

	provinceNum := cafe.ContactInfo.Province
	cityNum := cafe.ContactInfo.City

	publicCafe := PublicCafeProvinceCity{
		ID:           cafe.ID,
		OwnerID:      cafe.OwnerID,
		Name:         cafe.Name,
		Description:  cafe.Description,
		OpeningTime:  cafe.OpeningTime,
		ClosingTime:  cafe.ClosingTime,
		Comments:     cafe.Comments,
		Rating:       cafe.Rating,
		Images:       cafe.Images,
		Events:       cafe.Events,
		Capacity:     cafe.Capacity,
		ContactInfo:  cafe.ContactInfo,
		Categories:   cafe.Categories,
		ProvinceName: models.Provinces[provinceNum].Name,
		CityName:     models.Cities[provinceNum][cityNum].Name,
	}

	return &publicCafe, nil
}

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
	UserLastName  string `json:"user_last_name"`
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
			Comment:       comment,
			UserFirstName: userName.FirstName,
			UserLastName:  userName.LastName,
		}

		commentsWithNames = append(commentsWithNames, commentWithUserName)
	}

	return commentsWithNames, err
}

func (c CafeHandler) CreateEvent(ctx context.Context, event models.Event) error {
	start_time := event.StartTime.UTC()
	end_time := event.EndTime.UTC()

	if !utils.CheckStartTime(start_time) {
		return errors.ErrStartTimeInvalid.Error()
	}

	if !utils.CheckEndTime(start_time, end_time) {
		return errors.ErrEndTimeInvalid.Error()
	}

	eventID := rand.Int31()
	newEvent := &models.Event{
		ID:          eventID,
		CafeID:      event.CafeID,
		Name:        event.Name,
		Description: event.Description,
		StartTime:   start_time,
		EndTime:     end_time,
		ImageID:     event.ImageID,
	}

	if event.ImageID != "" {
		err := c.ImageRepo.Create(ctx, &models.Image{
			ID: event.ImageID,
			CafeID: event.CafeID,
		})

		if err != nil {
			log.GetLog().Errorf("Unable to create image. error: %v", err)
			return err
		}
	}

	err := c.EventRepo.CreateEventForCafe(ctx, newEvent)
	if err != nil {
		log.GetLog().Errorf("Unable to create event for cafe. error: %v", err)
		return err
	}

	return err
}

func (c CafeHandler) AddMenuItem(ctx context.Context, menuItem *models.MenuItem) (*models.MenuItem, error) {
	validCategory := false
	for _, category := range []models.MenuItemCategory{
		models.MenuItemCategoryCoffee,
		models.MenuItemCategoryTea,
		models.MenuItemCategoryAppetizer,
		models.MenuItemCategoryMainDish,
		models.MenuItemCategoryDessert,
		models.MenuItemCategoryDrink,
	} {
		if menuItem.Category == category {
			validCategory = true
			break
		}
	}

	if !validCategory {
		log.GetLog().Errorf("Category is invalid")
		return nil, fmt.Errorf("invalid menu item category: %s", menuItem.Category)
	}

	if menuItem.ImageID != "" {
		err := c.ImageRepo.Create(ctx, &models.Image{
			ID:     menuItem.ImageID,
			CafeID: menuItem.CafeID,
		})

		if err != nil {
			log.GetLog().Errorf("Unable to create image. error: %v", err)
			return nil, err
		}
	}

	err := c.MenuItemRepo.Create(ctx, menuItem)
	if err != nil {
		log.GetLog().Errorf("Unable to create menu item. error: %v", err)
		return nil, err
	}

	return menuItem, err
}

func (c CafeHandler) GetMenu(ctx context.Context, cafeID int32) ([]string, map[string][]*models.MenuItem, error) {
	menuItems, err := c.MenuItemRepo.GetItemsByCafeID(ctx, cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get menu items. error: %v", err)
		return nil, nil, err
	}

	var categories []string
	menu := make(map[string][]*models.MenuItem)

	for _, item := range menuItems {
		category := string(item.Category)
		categories = utils.AppendIfNotExists(categories, category)
		menu[category] = append(menu[category], item)
	}

	return categories, menu, nil
}
