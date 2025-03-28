package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/AaryanO2/go_projects/project_10_url_shortener/database"
	"github.com/AaryanO2/go_projects/project_10_url_shortener/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type request struct {
	URL         string        `json:"url"`
	Expiry      time.Duration `json:"expiry"`
	CustomShort string        `json:"short"`
}

type response struct {
	URL            string        `json:"url"`
	CustomShort    string        `json:"short"`
	Expiry         time.Duration `json:"expiry"`
	XRateRemaining int           `json:"rate_limit"`
	XRateLimitRest time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c fiber.Ctx) error {
	body := new(request)
	if err := c.Bind().JSON(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse body"})
	}

	// Rate limiting
	r2 := database.CreateClient(1)
	defer r2.Close()
	value, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*time.Minute)
		value = os.Getenv("API_QUOTA")
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Error connecting to database while rate limiting",
		})
	}

	valInt, _ := strconv.Atoi(value)
	if valInt < 1 {
		timeToReset, _ := r2.TTL(database.Ctx, c.IP()).Result()
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"ERROR":          "Rate limit exceeded",
			"LIMIT_RESETS_IN": timeToReset / time.Nanosecond / time.Minute ,
		})
	}

	// check input
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid url"})
	}

	//check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"Error": "Invalid url"})
	}
	//enforce http
	body.URL = helpers.EnforceHTTP(body.URL)

	// chec if short is already in use
	var id string

	r1 := database.CreateClient(0)
	defer r1.Close()

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
		val, err := r1.Get(database.Ctx, id).Result()
		if err != redis.Nil && err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Error": "Error connecting to the database while enforcing ssl",
			})
		} else if val != "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"Error": "URL already in use",
			})
		}
	}

	// Add to database

	if body.Expiry == 0 {
		body.Expiry = 24 * time.Hour
	}

	err = r1.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Unable to connect to server",
		})
	}
	resp := response{
		URL:            body.URL,
		CustomShort:    os.Getenv("DOMAIN") + "/" + id,
		Expiry:         body.Expiry,
		XRateRemaining: valInt - 1,
		XRateLimitRest: 30,
	}
	r2.Decr(database.Ctx, c.IP())

	resp.XRateLimitRest, _ = r2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitRest = resp.XRateLimitRest / time.Nanosecond / time.Minute

	return c.Status(fiber.StatusOK).JSON(resp)

}
