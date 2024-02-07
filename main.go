package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/akhil/go-fiber-postgres/models"
	"github.com/akhil/go-fiber-postgres/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type URLS struct {
	LongUrl  string `json:"longUrl"`
	ShortUrl string `json:"shortUrl"`
}
type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book has been added"})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(bookModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book delete successfully",
	})
	return nil
}

func (r *Repository) GetAllUrls(context *fiber.Ctx) error {
	urlModels := &[]models.URLS{}

	err := r.DB.Find(urlModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get urls something wen wrong"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "url's fetched successfully",
		"data":    urlModels,
	})
	return nil
}
func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get books"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books fetched successfully",
		"data":    bookModels,
	})
	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {

	id := context.Params("id")
	bookModel := &models.Books{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book id fetched successfully",
		"data":    bookModel,
	})
	return nil
}

func (r *Repository) GetUrl(context *fiber.Ctx) error {

	searchUrlModel := URLS{}
	result := URLS{}

	err := context.BodyParser(&searchUrlModel)

	if searchUrlModel.LongUrl == "" {
		if searchUrlModel.ShortUrl == "" {
			context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"message": "search url cannot be empty",
			})
			return nil
		}
		err = r.DB.Where("short_url = ?", searchUrlModel.ShortUrl).First(&result).Error
	} else {
		return nil
	}
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Url not found", "found": false})
		return err
	}
	// err=searchUrlModel.ShortUrl;
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Url fetched successfully",
		"url":     result,
	})
	return nil
}

func (r *Repository) GetShortUrl(context *fiber.Ctx, shortUrl string) error {

	searchUrlModel := URLS{}

	err := r.DB.Where("short_url = ?", shortUrl).First(&searchUrlModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Url not found", "found": false})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book id fetched successfully",
		"data":    searchUrlModel,
	})
	return nil
}

func (r *Repository) PostUrl(context *fiber.Ctx) error {
	url := URLS{}

	err := context.BodyParser(&url)

	shortUrl := randstr.Hex(16)
	getUrlErr := r.DB.Where("long_url = ?", url.LongUrl).First(&url).Error
	if getUrlErr != nil {
		url.ShortUrl = shortUrl
		err = r.DB.Create(&url).Error
		if err != nil {
			context.Status(http.StatusBadRequest).JSON(
				&fiber.Map{"message": "could not add Url"})
			return err
		}

		context.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "url has been added successfully", "url": url})
	} else {
		context.Status(http.StatusOK).JSON(
			&fiber.Map{"message": "this url is already added", "url": url})
	}

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, https://gofiber.net",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	api.Post("/create_books", r.CreateBook)
	api.Delete("delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
	api.Post("/addurl", r.PostUrl)
	api.Post("/geturl", r.GetUrl)
	api.Get("/get-all-urls", r.GetAllUrls)
	// app.Get("/:id", r.GetUrl)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}
	err = models.MigrateBooks(db)
	err = models.MigrateURL(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
