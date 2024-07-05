package models

type CafeCategory string

const (
	CafeCategoryCoffeeShop  CafeCategory = "coffee_shop"
	CafeCategoryRestaurant  CafeCategory = "restaurant"
	CafeCategoryBar         CafeCategory = "bar"
	CafeCategoryPub         CafeCategory = "pub"
	CafeCategoryBakery      CafeCategory = "bakery"
	CafeCategoryCafe        CafeCategory = "cafe"
	CafeCategoryTeaHouse    CafeCategory = "tea_house"
	CafeCategoryFastFood    CafeCategory = "fast_food"
	CafeCategoryFoodCourt   CafeCategory = "food_court"
	CafeCategoryDessertShop CafeCategory = "dessert_shop"
	CafeCategoryIceCream    CafeCategory = "ice_cream"
)

var CafeCategoryPersians = map[CafeCategory]string{
	CafeCategoryCoffeeShop:  "کافه",
	CafeCategoryRestaurant:  "رستوران",
	CafeCategoryBar:         "بار",
	CafeCategoryPub:         "پاب",
	CafeCategoryBakery:      "نانوایی",
	CafeCategoryCafe:        "کافه",
	CafeCategoryTeaHouse:    "چایخانه",
	CafeCategoryFastFood:    "فست فود",
	CafeCategoryFoodCourt:   "فودکورت",
	CafeCategoryDessertShop: "شیرینی",
	CafeCategoryIceCream:    "بستنی",
}

type MenuItemCategory string

const (
	MenuItemCategoryCoffee    MenuItemCategory = "coffee"
	MenuItemCategoryTea       MenuItemCategory = "tea"
	MenuItemCategoryAppetizer MenuItemCategory = "appetizer"
	MenuItemCategoryMainDish  MenuItemCategory = "main_dish"
	MenuItemCategoryDessert   MenuItemCategory = "dessert"
	MenuItemCategoryDrink     MenuItemCategory = "drink"
)

var MenuItemCategoryPersians = map[MenuItemCategory]string{
	MenuItemCategoryCoffee:    "قهوه",
	MenuItemCategoryTea:       "چای",
	MenuItemCategoryAppetizer: "پیش غذا",
	MenuItemCategoryMainDish:  "غذای اصلی",
	MenuItemCategoryDessert:   "دسر",
	MenuItemCategoryDrink:     "نوشیدنی",
}

type AmenityCategory string

const (
	AmenityCategoryWiFi              AmenityCategory = "wifi"
	AmenityCategoryTV                AmenityCategory = "tv"
	AmenityCategoryOutdoorSeating    AmenityCategory = "outdoor_seating"
	AmenityCategoryParking           AmenityCategory = "parking"
	AmenityCategoryLiveMusic         AmenityCategory = "live_music"
	AmenityCategoryBoardGames        AmenityCategory = "board_games"
	AmenityCategoryPetFriendly       AmenityCategory = "pet_friendly"
	AmenityCategoryHookah            AmenityCategory = "hookah"
	AmenityCategoryFreeSmoke         AmenityCategory = "free_smoke"
	AmenityCategoryVegetarianOptions AmenityCategory = "vegetarian_options"
	AmenityCategoryVeganOptions      AmenityCategory = "vegan_options"
	AmenityCategoryWheelchairAccess  AmenityCategory = "wheelchair_access"
	AmenityCategoryMeetingRooms      AmenityCategory = "meeting_rooms"
	AmenityCategoryWorkingSpace      AmenityCategory = "working_space"
)

var AmenityCategoryPersians = map[AmenityCategory]string{
	AmenityCategoryWiFi:              "وای فای",
	AmenityCategoryTV:                "تلویزیون",
	AmenityCategoryOutdoorSeating:    "فضای نشستن بیرون",
	AmenityCategoryParking:           "پارکینگ",
	AmenityCategoryLiveMusic:         "موسیقی زنده",
	AmenityCategoryBoardGames:        "بردگیم",
	AmenityCategoryPetFriendly:       "اجازه ورود حیوانات خانگی",
	AmenityCategoryHookah:            "قلیان",
	AmenityCategoryFreeSmoke:         "دود آزاد",
	AmenityCategoryVegetarianOptions: "غذای گیاهی",
	AmenityCategoryVeganOptions:      "غذای وگان",
	AmenityCategoryWheelchairAccess:  "دسترسی برای معلولان",
	AmenityCategoryMeetingRooms:      "اتاق جلسه",
	AmenityCategoryWorkingSpace:      "فضای کاری",
}
