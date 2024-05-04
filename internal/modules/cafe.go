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
			ID:        photoID,
			Reference: cafeID,
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
		images, err := c.ImageRepo.GetByReferenceID(ctx, cafes[i].ID)
		if err != nil {
			log.GetLog().Errorf("Unable to get cafe images. error: %v", err)
			continue
		}

		for _, image := range images {
			cafes[i].Images = append(cafes[i].Images, image.ID)
		}
	}

	return cafes, err
}

type PublicCafeProvinceCity struct {
	ID           int32                    `json:"id"`
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	OpeningTime  int8                     `json:"opening_time"`
	ClosingTime  int8                     `json:"closing_time"`
	Comments     []CommentWithUserName     `json:"comments"`
	Rating       float64                  `json:"rating"`
	Images       []string                 `json:"photos"`
	Events       []models.Event           `json:"events"`
	Capacity     int32                    `json:"capacity"`
	ContactInfo  models.ContactInfo       `json:"contact_info"`
	Categories   []models.CafeCategory    `json:"categories"`
	Amenities    []models.AmenityCategory `json:"amenities"`
	ProvinceName string                   `json:"province_name"`
	CityName     string                   `json:"city_name"`
}

func (c CafeHandler) PublicCafeProfile(ctx context.Context, cafeID int32) (*PublicCafeProvinceCity, error) {
	cafe, err := c.CafeRepo.GetByID(ctx, int32(cafeID))
	if err != nil {
		log.GetLog().Errorf("Cafe id does not exist. error: %v", err)
		return nil, err
	}

	comments, err := c.CommentRepo.GetAllByCafeID(ctx, int32(cafeID))
	if err != nil {
		log.GetLog().Errorf("Unable to get all comments. error: %v", err)
		return nil, err
	}

	commentsWithNames := make([]CommentWithUserName, len(comments))
	for i, comment := range comments {
		userName, err := c.UserRepo.GetByID(ctx, comment.UserID)
		if err != nil {
			log.GetLog().Errorf("Unable to get user name for user ID %d: %v", comment.UserID, err)
			return nil, err
		}
		
		commentsWithNames[i].Comment = comment
		commentsWithNames[i].UserFirstName = userName.FirstName
		commentsWithNames[i].UserLastName = userName.LastName
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

	photos, err := c.ImageRepo.GetByReferenceID(ctx, int32(cafeID))
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
		Name:         cafe.Name,
		Description:  cafe.Description,
		OpeningTime:  cafe.OpeningTime,
		ClosingTime:  cafe.ClosingTime,
		Comments:     commentsWithNames,
		Rating:       cafe.Rating,
		Images:       cafe.Images,
		Events:       cafe.Events,
		Capacity:     cafe.Capacity,
		ContactInfo:  cafe.ContactInfo,
		Categories:   cafe.Categories,
		Amenities:    cafe.Amenities,
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
	UserFirstName string `json:"first_name"`
	UserLastName  string `json:"last_name"`
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
			ID:        event.ImageID,
			Reference: eventID,
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

	itemID, err := c.MenuItemRepo.Create(ctx, menuItem)
	if err != nil {
		log.GetLog().Errorf("Unable to create menu item. error: %v", err)
		return nil, err
	}

	if menuItem.ImageID != "" {
		err := c.ImageRepo.Create(ctx, &models.Image{
			ID:        menuItem.ImageID,
			Reference: itemID,
		})

		if err != nil {
			log.GetLog().Errorf("Unable to create image. error: %v", err)
			return nil, err
		}
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

func (c CafeHandler) EditMenuItem(ctx context.Context, newItem models.MenuItem) error {
	preItem, err := c.MenuItemRepo.GetByID(ctx, newItem.ID)
	if err != nil {
		log.GetLog().Errorf("Incorrect menu item id. error: %v", err)
		return err
	}

	images, err := c.ImageRepo.GetByReferenceID(ctx, preItem.ID)
	if err != nil {
		log.GetLog().Errorf("Unable to get image by reference id. error: %v", err)
		return err
	}
	
	if len(images) != 0 {
		preItem.ImageID = images[0].ID
	}

	if newItem.Name != preItem.Name && newItem.Name != "" {
		err = c.MenuItemRepo.UpdateName(ctx, newItem.ID, newItem.Name)
		if err != nil {
			log.GetLog().Errorf("Unable to update menu items name. error: %v", err)
			return err
		}
	}

	if newItem.Price != preItem.Price && newItem.Price != 0 {
		err = c.MenuItemRepo.UpdatePrice(ctx, newItem.ID, newItem.Price)
		if err != nil {
			log.GetLog().Errorf("Unable to update menu items price. error: %v", err)
			return err
		}
	}

	if len(newItem.Ingredients) != len(preItem.Ingredients) && len(newItem.Ingredients) != 0 {
		err = c.MenuItemRepo.UpdateIngredients(ctx, newItem.ID, newItem.Ingredients)
		if err != nil {
			log.GetLog().Errorf("Unable to update menu items ingredients. error: %v", err)
			return err
		}
	}

	if newItem.ImageID != preItem.ImageID && newItem.ImageID != "" {
		err = c.MenuItemRepo.UpdateImageID(ctx, newItem.ID, newItem.ImageID)
		if err != nil {
			log.GetLog().Errorf("Unable to update menu items image. error: %v", err)
			return err
		}
	}

	return err
}

func (c CafeHandler) DeleteMenuItem(ctx context.Context, itemID int32) error {
	err := c.MenuItemRepo.DeleteByID(ctx, itemID)
	if err != nil {
		log.GetLog().Errorf("Unable to delete item by menu item id. error: %v", err)
		return err
	}

	err = c.ImageRepo.DeleteByReferenceID(ctx, itemID)
	if err != nil {
		log.GetLog().Errorf("Unable to delete image by menu item id. error: %v", err)
		return err
	}

	return err
}

func (c CafeHandler) Home(ctx context.Context) ([]models.Cafe, []*models.Comment, error) {
	cafes, err := c.Rating.GetNTopRatings(ctx, 5)
	if err != nil {
		log.GetLog().Errorf("Unable to get home cafes. error: %v", err)
		return nil, nil, err
	}

	cafes = append(cafes, 1)
	ds, err := c.CafeRepo.GetByCafeIDs(ctx, cafes)
	if err != nil {
		return nil, nil, err
	}

	comments, err := c.CommentRepo.GetLast(ctx, 0, 5, 0)
	if err != nil {
		log.GetLog().Errorf("Unable to get home comments. error: %v", err)
		return nil, nil, err
	}

	return ds, comments, nil

}
