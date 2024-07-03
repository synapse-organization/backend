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

	"github.com/redis/go-redis/v9"
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
	PaymentRepo     repo.Transaction
	LocationsRepo   repo.LocationsRepo
	FavoriteRepo    repo.FavoritesRepo
	Redis           *redis.Client
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
		cafes[i].Location, err = c.LocationsRepo.GetCafeLocation(ctx, cafe.ID)
		if err != nil {
			log.GetLog().Errorf("Unable to get cafe rating. error: %v", err)
		}
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
	ID               int32                    `json:"id"`
	Name             string                   `json:"name"`
	Description      string                   `json:"description"`
	OpeningTime      int8                     `json:"opening_time"`
	ClosingTime      int8                     `json:"closing_time"`
	Comments         []CommentWithUserName    `json:"comments"`
	Rating           float64                  `json:"rating"`
	Images           []string                 `json:"photos"`
	Events           []models.Event           `json:"events"`
	Capacity         int32                    `json:"capacity"`
	ContactInfo      models.ContactInfo       `json:"contact_info"`
	Categories       []models.CafeCategory    `json:"categories"`
	Amenities        []models.AmenityCategory `json:"amenities"`
	ProvinceName     string                   `json:"province_name"`
	CityName         string                   `json:"city_name"`
	ReservationPrice float64                  `json:"reservation_price"`
	Favorite         bool                     `json:"favorite"`
}

func (c CafeHandler) PublicCafeProfile(ctx context.Context, cafeID int32, userID int32) (*PublicCafeProvinceCity, error) {
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
		images, err := c.ImageRepo.GetByReferenceID(ctx, event.ID)
		if err != nil {
			log.GetLog().Errorf("Unable to get images by event id. error: %v", err)
			return nil, err
		}

		if images != nil {
			cafe.Events[i].ImageID = images[0].ID
		}
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

	var isFavorite bool
	if userID != 0 {
		isFavorite, err = c.FavoriteRepo.CheckExists(ctx, userID, cafeID)
		if err != nil {
			log.GetLog().Errorf("Unable to check favorite existence. error: %v", err)
			return nil, err
		}
	} else {
		isFavorite = false
	}

	publicCafe := PublicCafeProvinceCity{
		ID:               cafe.ID,
		Name:             cafe.Name,
		Description:      cafe.Description,
		OpeningTime:      cafe.OpeningTime,
		ClosingTime:      cafe.ClosingTime,
		Comments:         commentsWithNames,
		Rating:           cafe.Rating,
		Images:           cafe.Images,
		Events:           cafe.Events,
		Capacity:         cafe.Capacity,
		ContactInfo:      cafe.ContactInfo,
		Categories:       cafe.Categories,
		Amenities:        cafe.Amenities,
		ProvinceName:     models.Provinces[provinceNum-1].Name,
		CityName:         models.Cities[cityNum-1].Name,
		ReservationPrice: cafe.ReservationPrice,
		Favorite:         isFavorite,
	}

	return &publicCafe, nil
}

type PrivateCafeProvinceCity struct {
	ID               int32                    `json:"id"`
	Name             string                   `json:"name"`
	Description      string                   `json:"description"`
	OpeningTime      int8                     `json:"opening_time"`
	ClosingTime      int8                     `json:"closing_time"`
	Comments         []CommentWithUserName    `json:"comments"`
	Rating           float64                  `json:"rating"`
	Images           []string                 `json:"photos"`
	Events           []models.Event           `json:"events"`
	Capacity         int32                    `json:"capacity"`
	ContactInfo      models.ContactInfo       `json:"contact_info"`
	Categories       []models.CafeCategory    `json:"categories"`
	Amenities        []models.AmenityCategory `json:"amenities"`
	ProvinceName     string                   `json:"province_name"`
	CityName         string                   `json:"city_name"`
	ReservationPrice float64                  `json:"reservation_price"`
}

func (c CafeHandler) PrivateCafeProfile(ctx context.Context, cafeID int32) (*PrivateCafeProvinceCity, error) {
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
		images, err := c.ImageRepo.GetByReferenceID(ctx, event.ID)
		if err != nil {
			log.GetLog().Errorf("Unable to get images by event id. error: %v", err)
			return nil, err
		}

		if images != nil {
			cafe.Events[i].ImageID = images[0].ID
		}
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

	privateCafe := PrivateCafeProvinceCity{
		ID:               cafe.ID,
		Name:             cafe.Name,
		Description:      cafe.Description,
		OpeningTime:      cafe.OpeningTime,
		ClosingTime:      cafe.ClosingTime,
		Comments:         commentsWithNames,
		Rating:           cafe.Rating,
		Images:           cafe.Images,
		Events:           cafe.Events,
		Capacity:         cafe.Capacity,
		ContactInfo:      cafe.ContactInfo,
		Categories:       cafe.Categories,
		Amenities:        cafe.Amenities,
		ProvinceName:     models.Provinces[provinceNum-1].Name,
		CityName:         models.Cities[cityNum-1].Name,
		ReservationPrice: cafe.ReservationPrice,
	}

	return &privateCafe, nil
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

	if !utils.CheckPriceValidity(event.Price) {
		return errors.ErrPriceInvalid.Error()
	}

	if !utils.CheckCapacityValidity(event.Capacity) {
		return errors.ErrCapacityInvalid.Error()
	}

	eventID := rand.Int31()
	newEvent := &models.Event{
		ID:               eventID,
		CafeID:           event.CafeID,
		Name:             event.Name,
		Description:      event.Description,
		StartTime:        start_time,
		EndTime:          end_time,
		ImageID:          event.ImageID,
		Price:            event.Price,
		Capacity:         event.Capacity,
		CurrentAttendees: 0,
		Reservable:       true,
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

func (c CafeHandler) GetMenu(ctx context.Context, cafeID int32) ([]string, map[string][]*models.MenuItem, string, string, error) {
	menuItems, err := c.MenuItemRepo.GetItemsByCafeID(ctx, cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get menu items. error: %v", err)
		return nil, nil, "", "", err
	}

	for i, item := range menuItems {
		images, err := c.ImageRepo.GetByReferenceID(ctx, item.ID)
		if err != nil {
			log.GetLog().Errorf("Unable to get image by reference id. error: %v", err)
			return nil, nil, "", "", err
		}

		if len(images) > 0 {
			menuItems[i].ImageID = images[0].ID
		}
	}

	cafe, err := c.CafeRepo.GetByID(ctx, cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get cafe. error: %v", err)
		return nil, nil, "", "", err
	}

	imageID, err := c.ImageRepo.GetMainImage(ctx, cafe.ID)
	if err != nil {
		log.GetLog().Errorf("Unable to get cafe. error: %v", err)
		return nil, nil, "", "", err
	}

	var categories []string
	menu := make(map[string][]*models.MenuItem)

	for _, item := range menuItems {
		category := string(item.Category)
		categories = utils.AppendIfNotExists(categories, category)
		menu[category] = append(menu[category], item)
	}

	return categories, menu, cafe.Name, imageID, nil
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

	if newItem.ImageID != "" {
		err := c.ImageRepo.DeleteByID(ctx, newItem.ImageID)
		if err != nil {
			log.GetLog().Errorf("Unable to delete image by id. error: %v", err)
			return err
		}

		if newItem.ImageID != preItem.ImageID {
			imageID := string(rand.Int31())
			err = c.ImageRepo.Create(ctx, &models.Image{
				ID:        imageID,
				Reference: newItem.ID,
			})
			if err != nil {
				log.GetLog().Errorf("Unable to create image. error: %v", err)
				return err
			}
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

func (c CafeHandler) Home(ctx context.Context) ([]models.Cafe, []*models.Comment, []*models.Event, error) {
	cafes, err := c.Rating.GetNTopRatings(ctx, 5)
	if err != nil {
		log.GetLog().Errorf("Unable to get home cafes. error: %v", err)
		return nil, nil, nil, err
	}

	cafes = append(cafes, 1)
	ds, err := c.CafeRepo.GetByCafeIDs(ctx, cafes)
	if err != nil {
		return nil, nil, nil, err
	}

	comments, err := c.CommentRepo.GetLast(ctx, 0, 5, 0)
	if err != nil {
		log.GetLog().Errorf("Unable to get home comments. error: %v", err)
		return nil, nil, nil, err
	}

	event, err := c.EventRepo.GetAllEventsNearestStartTime(ctx, 5)
	if err != nil {
		log.GetLog().Errorf("Unable to get home events. error: %v", err)
		return nil, nil, nil, err
	}

	return ds, comments, event, nil
}

func (c CafeHandler) ReserveEvent(ctx context.Context, eventID int32, userID int32) error {
	event, err := c.EventRepo.GetEventByID(ctx, eventID)
	if err != nil {
		log.GetLog().Errorf("Unable to get event by id. error: %v", err)
		return err
	}

	cafe, err := c.CafeRepo.GetByID(ctx, event.CafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get cafe by id. error: %v", err)
		return err
	}

	userEvents, err := c.EventRepo.GetEventsByUserID(ctx, userID)
	if err != nil {
		log.GetLog().Errorf("Unable to get event by user id. error: %v", err)
		return err
	}

	for _, event := range userEvents {
		if event.ID == eventID {
			return errors.ErrEventReserved.Error()
		}
	}

	if !event.Reservable {
		return errors.ErrEventUnreservable.Error()
	}

	transactionID, err := c.PaymentRepo.Create(ctx, &models.Transaction{
		SenderID:    userID,
		ReceiverID:  cafe.OwnerID,
		Amount:      int64(event.Price),
		Description: event.Description,
		Type:        3,
		CreatedAt:   time.Now().UTC(),
	})
	if err != nil {
		log.GetLog().Errorf("Unable to do transaction. error: %v", err)
		return err
	}

	err = c.EventRepo.CreateEventForUser(ctx, userID, eventID, transactionID)
	if err != nil {
		log.GetLog().Errorf("Unable to create event for user. error: %v", err)
		return err
	}

	event.CurrentAttendees++
	if event.CurrentAttendees == event.Capacity {
		event.Reservable = false
		err = c.EventRepo.UpdateEvent(ctx, eventID, repo.UpdateEventReservability, false)
		if err != nil {
			log.GetLog().Errorf("Unable to update event reservability by id. error: %v", err)
			return err
		}
	}
	err = c.EventRepo.UpdateEvent(ctx, eventID, repo.UpdateEventAttendees, event.CurrentAttendees)
	if err != nil {
		log.GetLog().Errorf("Unable to update event attendees by id. error: %v", err)
		return err
	}

	return nil
}

// type PrivateCafeRes struct {
// 	ID           int32                    `json:"id"`
// 	Name         string                   `json:"name"`
// 	Description  string                   `json:"description"`
// 	OpeningTime  int8                     `json:"opening_time"`
// 	ClosingTime  int8                     `json:"closing_time"`
// 	Comments     []CommentWithUserName    `json:"comments"`
// 	Rating       float64                  `json:"rating"`
// 	Images       []string                 `json:"photos"`
// 	Events       []models.Event           `json:"events"`
// 	Reservations []models.Reservation     `json:"reservations"`
// 	Capacity     int32                    `json:"capacity"`
// 	ContactInfo  models.ContactInfo       `json:"contact_info"`
// 	Categories   []models.CafeCategory    `json:"categories"`
// 	Amenities    []models.AmenityCategory `json:"amenities"`
// 	ProvinceName string                   `json:"province_name"`
// 	CityName     string                   `json:"city_name"`
// }

// func (c CafeHandler) PrivateCafe(ctx context.Context, cafe models.Cafe) (*PrivateCafeRes, error) {
// 	publicCafe, err := c.PublicCafeProfile(ctx, cafe.ID)
// 	if err != nil {
// 		log.GetLog().Errorf("Unable to get public cafe info. error: %v", err)
// 		return nil, err
// 	}

// 	reservations, err := c.ReservationRepo.GetByCafeID(ctx, cafe.ID)
// 	if err != nil {
// 		log.GetLog().Errorf("Unable to get reservations by cafe id. error: %v", err)
// 		return nil, err
// 	}

// 	cafe.Reservations = make([]models.Reservation, len(reservations))
// 	for i, reservation := range reservations {
// 		cafe.Reservations[i] = *reservation
// 	}

// 	privateCafe := PrivateCafeRes{
// 		ID:           publicCafe.ID,
// 		Name:         publicCafe.Name,
// 		Description:  publicCafe.Description,
// 		OpeningTime:  publicCafe.OpeningTime,
// 		ClosingTime:  publicCafe.ClosingTime,
// 		Comments:     publicCafe.Comments,
// 		Rating:       publicCafe.Rating,
// 		Images:       publicCafe.Images,
// 		Events:       publicCafe.Events,
// 		Reservations: cafe.Reservations,
// 		Capacity:     publicCafe.Capacity,
// 		ContactInfo:  publicCafe.ContactInfo,
// 		Categories:   publicCafe.Categories,
// 		Amenities:    publicCafe.Amenities,
// 		ProvinceName: publicCafe.ProvinceName,
// 		CityName:     publicCafe.CityName,
// 	}

// 	return &privateCafe, nil
// }

type RequestEditCafe struct {
	ID               int32                    `json:"id"`
	Name             string                   `json:"name"`
	Description      string                   `json:"description"`
	OpeningTime      int8                     `json:"opening_time"`
	ClosingTime      int8                     `json:"closing_time"`
	DeletedImages    []string                 `json:"deleted_images"`
	AddedImages      []string                 `json:"added_images"`
	Capacity         int32                    `json:"capacity"`
	ContactInfo      models.ContactInfo       `json:"contact_info"`
	Categories       []models.CafeCategory    `json:"categories"`
	Amenities        []models.AmenityCategory `json:"amenities"`
	ReservationPrice float64                  `json:"reservation_price"`
}

func (c CafeHandler) EditCafe(ctx context.Context, newCafe RequestEditCafe) error {
	preCafe, err := c.CafeRepo.GetByID(ctx, newCafe.ID)
	if err != nil {
		log.GetLog().Errorf("Incorrect cafe id. error: %v", err)
		return err
	}

	images, err := c.ImageRepo.GetByReferenceID(ctx, newCafe.ID)
	if err != nil {
		log.GetLog().Errorf("Unable to get image by reference id. error: %v", err)
		return err
	}

	preCafe.Images = make([]string, len(images))
	for i, image := range images {
		preCafe.Images[i] = image.ID
	}

	if preCafe.Name != newCafe.Name {
		err = c.CafeRepo.Update(ctx, newCafe.ID, repo.UpdateCafeName, newCafe.Name)
		if err != nil {
			log.GetLog().Errorf("Unable to update cafes name. error: %v", err)
			return err
		}
	}

	if preCafe.Description != newCafe.Description {
		err = c.CafeRepo.Update(ctx, newCafe.ID, repo.UpdateCafeDescription, newCafe.Description)
		if err != nil {
			log.GetLog().Errorf("Unable to update cafes description. error: %v", err)
			return err
		}
	}

	if preCafe.OpeningTime != newCafe.OpeningTime {
		if !utils.CheckCafeTimeValidity(newCafe.OpeningTime) {
			errors.ErrStartTimeInvalid.Error()
		}

		err = c.CafeRepo.Update(ctx, newCafe.ID, repo.UpdateCafeOpeningTime, newCafe.OpeningTime)
		if err != nil {
			log.GetLog().Errorf("Unable to update cafes opening time. error: %v", err)
			return err
		}
	}

	if preCafe.ClosingTime != newCafe.ClosingTime {
		if !utils.CheckCafeTimeValidity(newCafe.ClosingTime) {
			errors.ErrEndTimeInvalid.Error()
		}

		err = c.CafeRepo.Update(ctx, newCafe.ID, repo.UpdateCafeClosingTime, newCafe.ClosingTime)
		if err != nil {
			log.GetLog().Errorf("Unable to update cafes closing time. error: %v", err)
			return err
		}
	}

	if preCafe.Capacity != newCafe.Capacity {
		if !utils.CheckCapacityValidity(newCafe.Capacity) {
			return errors.ErrCapacityInvalid.Error()
		}

		err = c.CafeRepo.Update(ctx, newCafe.ID, repo.UpdateCafeCapacity, newCafe.Capacity)
		if err != nil {
			log.GetLog().Errorf("Unable to update cafes capacity. error: %v", err)
			return err
		}
	}

	if preCafe.ContactInfo.Phone != newCafe.ContactInfo.Phone {
		if !utils.CheckPhoneValidity(newCafe.ContactInfo.Phone) {
			return errors.ErrPhoneInvalid.Error()
		}

		err = c.CafeRepo.Update(ctx, newCafe.ID, repo.UpdateCafePhoneNumber, newCafe.ContactInfo.Phone)
		if err != nil {
			log.GetLog().Errorf("Unable to update cafes phone number. error: %v", err)
			return err
		}
	}

	if preCafe.ContactInfo.Email != newCafe.ContactInfo.Email {
		if !utils.CheckEmailValidity(newCafe.ContactInfo.Email) {
			return errors.ErrEmailExists.Error()
		}

		err = c.CafeRepo.Update(ctx, newCafe.ID, repo.UpdateCafeEmail, newCafe.ContactInfo.Email)
		if err != nil {
			log.GetLog().Errorf("Unable to update cafes email. error: %v", err)
			return err
		}
	}

	if preCafe.ContactInfo.Province != newCafe.ContactInfo.Province {
		err = c.CafeRepo.Update(ctx, newCafe.ID, repo.UpdateCafeProvince, newCafe.ContactInfo.Province)
		if err != nil {
			log.GetLog().Errorf("Unable to update cafes province. error: %v", err)
			return err
		}
	}

	if preCafe.ContactInfo.City != newCafe.ContactInfo.City {
		err = c.CafeRepo.Update(ctx, newCafe.ID, repo.UpdateCafeCity, newCafe.ContactInfo.City)
		if err != nil {
			log.GetLog().Errorf("Unable to update cafes city. error: %v", err)
			return err
		}
	}

	if preCafe.ContactInfo.Address != newCafe.ContactInfo.Address {
		err = c.CafeRepo.Update(ctx, newCafe.ID, repo.UpdateCafeAddress, newCafe.ContactInfo.Address)
		if err != nil {
			log.GetLog().Errorf("Unable to update cafes address. error: %v", err)
			return err
		}
	}

	err = c.CafeRepo.Update(ctx, newCafe.ID, repo.UpdateCafeCategories, newCafe.Categories)
	if err != nil {
		log.GetLog().Errorf("Unable to update cafes categories. error: %v", err)
		return err
	}

	err = c.CafeRepo.Update(ctx, newCafe.ID, repo.UpdateCafeAmenities, newCafe.Amenities)
	if err != nil {
		log.GetLog().Errorf("Unable to update cafes amenities. error: %v", err)
		return err
	}

	for _, imageID := range newCafe.DeletedImages {
		err = c.ImageRepo.DeleteByID(ctx, imageID)
		if err != nil {
			log.GetLog().Errorf("Unable to delete image by id. error: %v", err)
			return err
		}
	}

	for _, imageID := range newCafe.AddedImages {
		if imageID != "" {
			err = c.ImageRepo.Create(ctx, &models.Image{
				ID:        imageID,
				Reference: newCafe.ID,
			})
			if err != nil {
				log.GetLog().Errorf("Unable to add image by id. error: %v", err)
				return err
			}
		}
	}

	if preCafe.ReservationPrice != newCafe.ReservationPrice {
		err = c.CafeRepo.Update(ctx, newCafe.ID, repo.UpdateCafeReservationPrice, newCafe.ReservationPrice)
		if err != nil {
			log.GetLog().Errorf("Unable to update cafe reservation price. error: %v", err)
			return err
		}
	}

	return err
}

func (c CafeHandler) EditEvent(ctx context.Context, newEvent models.Event) error {
	preEvent, err := c.EventRepo.GetEventByID(ctx, newEvent.ID)
	if err != nil {
		log.GetLog().Errorf("Incorrect event id. error: %v", err)
		return err
	}

	images, err := c.ImageRepo.GetByReferenceID(ctx, newEvent.ID)
	if err != nil {
		log.GetLog().Errorf("Unable to get image by reference id. error: %v", err)
		return err
	}

	if images != nil {
		preEvent.ImageID = images[0].ID
	}

	newFields := []interface{}{newEvent.Name, newEvent.Description, newEvent.Price}
	preFields := []interface{}{preEvent.Name, preEvent.Description, preEvent.Price}
	updateFields := []repo.UpdateEventType{repo.UpdateEventName, repo.UpdateEventDescription, repo.UpdateEventPrice}

	for i := range newFields {
		if preFields[i] != newFields[i] {
			err = c.EventRepo.UpdateEvent(ctx, newEvent.ID, updateFields[i], newFields[i])
			if err != nil {
				log.GetLog().Errorf("Unable to update event by id. error: %v", err)
				return err
			}
		}
	}

	if preEvent.Capacity != newEvent.Capacity {
		updateCapacity, updateReserve, err := utils.CheckReservability(preEvent.Reservable, newEvent.Reservable, preEvent.Capacity, newEvent.Capacity, preEvent.CurrentAttendees)
		if err != nil {
			log.GetLog().Errorf("Invalid new capacity. error: %v", err)
			return err
		}

		if !updateCapacity {
			err = c.EventRepo.UpdateEvent(ctx, newEvent.ID, repo.UpdateEventCapacity, newEvent.Capacity)
			if err != nil {
				log.GetLog().Errorf("Unable to update event by id. error: %v", err)
				return err
			}
		}

		if !updateReserve {
			err = c.EventRepo.UpdateEvent(ctx, newEvent.ID, repo.UpdateEventReservability, !preEvent.Reservable)
			if err != nil {
				log.GetLog().Errorf("Unable to update event by id. error: %v", err)
				return err
			}
		}
	}

	if preEvent.ImageID != newEvent.ImageID {
		if preEvent.ImageID != "" {
			err = c.ImageRepo.DeleteByID(ctx, preEvent.ImageID)
			if err != nil {
				log.GetLog().Errorf("Unable to delete image by id. error: %v", err)
				return err
			}
		}

		if newEvent.ImageID != "" {
			eventID := string(rand.Int31())
			err = c.ImageRepo.Create(ctx, &models.Image{
				ID:        eventID,
				Reference: newEvent.ID,
			})
			if err != nil {
				log.GetLog().Errorf("Unable to create image. error: %v", err)
				return err
			}
		}
	}

	return err
}

func (c CafeHandler) DeleteEvent(ctx context.Context, eventID int32) error {
	err := c.EventRepo.DeleteByID(ctx, eventID)
	if err != nil {
		log.GetLog().Errorf("Unable to delete event by id. error: %v", err)
		return err
	}

	err = c.ImageRepo.DeleteByReferenceID(ctx, eventID)
	if err != nil {
		log.GetLog().Errorf("Unable to delete image by reference id. error: %v", err)
		return err
	}

	return err
}

func (c CafeHandler) GetFullyBookedDays(ctx context.Context, cafeID int32, startDate time.Time) ([]time.Time, error) {
	cafe, err := c.CafeRepo.GetByID(ctx, cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get cafe. error: %v", err)
		return nil, err
	}

	return c.ReservationRepo.GetFullyBookedDays(ctx, cafeID, startDate, cafe.OpeningTime, cafe.ClosingTime)
}

func (c CafeHandler) GetAvailableTimeSlots(ctx context.Context, cafeID int32, day time.Time) ([]map[string]interface{}, error) {
	cafe, err := c.CafeRepo.GetByID(ctx, cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get cafe. error: %v", err)
		return nil, err
	}

	return c.ReservationRepo.GetAvailableTimeSlots(ctx, cafeID, day, cafe.Capacity, cafe.OpeningTime, cafe.ClosingTime)
}

func (c CafeHandler) ReserveCafe(ctx context.Context, reservation *models.Reservation) error {
	totalPeople, err := c.ReservationRepo.CountByTime(ctx, reservation.CafeID, reservation.StartTime, reservation.EndTime)
	if err != nil {
		log.GetLog().Errorf("Unable to check availability. error: %v", err)
		return err
	}

	cafe, err := c.CafeRepo.GetByID(ctx, reservation.CafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get cafe. error: %v", err)
		return err
	}

	if totalPeople+reservation.People > cafe.Capacity {
		log.GetLog().Errorf("time slot is fully booked")
		return fmt.Errorf("time slot is fully booked")
	}

	transactionID, err := c.PaymentRepo.Create(ctx, &models.Transaction{
		SenderID:    reservation.UserID,
		ReceiverID:  cafe.OwnerID,
		Amount:      int64(cafe.ReservationPrice * float64(reservation.People)),
		Description: "cafe reservation transaction",
		Type:        3,
		CreatedAt:   time.Now().UTC(),
	})
	if err != nil {
		log.GetLog().Errorf("Unable to do transaction. error: %v", err)
		return err
	}

	reservation.TransactionID = transactionID
	err = c.ReservationRepo.Create(ctx, reservation)
	if err != nil {
		log.GetLog().Errorf("Unable to create reservation. error: %v", err)
		return err
	}

	return nil
}

func (c CafeHandler) GetNearestCafes(ctx context.Context, lat float64, long float64, radius float64) ([]redis.GeoLocation, error) {

	return c.Redis.GeoRadius(ctx, "locations", lat, long, &redis.GeoRadiusQuery{
		Radius:      radius,
		Unit:        "km",
		WithCoord:   true,
		WithDist:    true,
		WithGeoHash: true,
		Count:       5,
		Sort:        "ASC",
	}).Result()

}

func (c CafeHandler) SetCafeLocation(ctx context.Context, m *models.Location) error {
	return c.LocationsRepo.SetLocation(ctx, m)
}

func (c CafeHandler) GetCafeLocation(ctx context.Context, id int32) (models.Location, error) {
	return c.LocationsRepo.GetCafeLocation(ctx, id)
}

func (c CafeHandler) AddRating(ctx context.Context, userID, cafeID, rating int32) error {
	return c.Rating.CreateOrUpdate(ctx, &models.Rating{
		ID:     rand.Int31(),
		UserID: userID,
		CafeID: cafeID,
		Rating: rating,
	})
}

type ReservationInfo struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	People    int32     `json:"people"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (c CafeHandler) GetCafeReservations(ctx context.Context, cafe *models.Cafe, day time.Time) ([]ReservationInfo, error) {
	reservations, err := c.ReservationRepo.GetByDateCafeID(ctx, cafe.ID, day, day.Add(time.Hour*24))
	if err != nil {
		log.GetLog().Errorf("Unable to get reservations by cafe id. error: %v", err)
		return nil, err
	}

	reservationsInfo := []ReservationInfo{}

	for _, reservation := range *reservations {
		user, err := c.UserRepo.GetByID(ctx, reservation.UserID)
		if err != nil {
			log.GetLog().Errorf("Unable to get user by id. error: %v", err)
			return nil, err
		}

		reservationsInfo = append(reservationsInfo, ReservationInfo{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			People:    reservation.People,
			StartTime: reservation.StartTime,
			EndTime:   reservation.EndTime,
		})
	}

	return reservationsInfo, nil
}

func (c CafeHandler) AddToFavorite(ctx context.Context, userID int32, cafeID int32) error {
	exists, err := c.FavoriteRepo.CheckExists(ctx, userID, cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to check favorite existence. error: %v", err)
		return err
	}

	if exists {
		return fmt.Errorf("Favorite already exists")
	}

	_, err = c.FavoriteRepo.Create(ctx, &models.Favorite{
		UserID: userID,
		CafeID: cafeID,
	})
	if err != nil {
		log.GetLog().Errorf("Unable to create favorite. error: %v", err)
		return err
	}

	return nil
}

func (c CafeHandler) GetRating(ctx context.Context, userID int32, cafeID int32) (float64, int32, error) {
	cafeRating, err := c.Rating.GetCafesRating(ctx, cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get cafe rating. error: %v", err)
		return 0, 0, err
	}

	cafeRatingCount, err := c.Rating.GetRatingByUserIDAndCafeID(ctx, userID, cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to get cafe rating count. error: %v", err)
		return 0, 0, err
	}

	return cafeRating, cafeRatingCount.Rating, nil
}

func (c CafeHandler) RemoveFavorite(ctx context.Context, userID int32, cafeID int32) error {
	err := c.FavoriteRepo.DeleteByIDs(ctx, userID, cafeID)
	if err != nil {
		log.GetLog().Errorf("Unable to delete favorite. error: %v", err)
		return err
	}

	return nil
}

func (c CafeHandler) GetFavoriteList(ctx context.Context, toInt32 int32) ([]models.Cafe, error) {

	favorites, err := c.FavoriteRepo.GetFavoritesByUserID(ctx, toInt32)
	if err != nil {
		log.GetLog().Errorf("Unable to get favorite list. error: %v", err)
		return nil, err
	}

	var cafes []models.Cafe
	for _, favorite := range favorites {
		cafe, err := c.CafeRepo.GetByID(ctx, favorite.CafeID)
		if err != nil {
			log.GetLog().Errorf("Unable to get cafe by id. error: %v", err)
			return nil, err
		}

		cafes = append(cafes, *cafe)
	}

	return cafes, nil

}
