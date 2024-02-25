package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/ShaunakSensarma/URL-Shortener/database"
	"github.com/ShaunakSensarma/URL-Shortener/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

/*
defining request and response structure to have a format the the frontend will expect.
makes sure the code is very stable.
*/
type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

// shortenURL is responsible for shortening the URL.
func ShortenURL(c *fiber.Ctx) error {
	body := new(request)

	// to parse the input JSON into struct defined above.
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"eror": "cannot parse JSON"})
	}

	//-----------------------------------------------------------------------------------------------
	// IMPLEMENT RATE LIMITER (RESETS EVERY 30 MINUTES)
	//-----------------------------------------------------------------------------------------------
	r2 := database.CreateClient(1)
	defer r2.Close()

	_, err := r2.Get(database.Ctx, c.IP()).Result() //getting value associated with the key: IP (more interested in error).
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		quota_val, _ := r2.Get(database.Ctx, c.IP()).Result() //we are getting remaining quota.
		valInt, _ := strconv.Atoi(quota_val)
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	//-----------------------------------------------------------------------------------------------

	// check if the input is an actual URL.
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	// check for domain error.
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "You have gone to bad address"})
	}

	//enforce https, SSL.
	body.URL = helpers.EnforceHTTP(body.URL)

	//----------------------------------------------------------------------------------------------
	// CREATING CUSTOM SHORT URL (ENSURE NOT ALREADY IN USE).
	//----------------------------------------------------------------------------------------------
	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	// checking whether the short_url is already in use.
	short_url_val, _ := r.Get(database.Ctx, id).Result()
	if short_url_val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL Custom Short is already in use",
		})
	}

	// re-setting expiry time of the new custom short url.
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	// setting the entire thing in database. (updating in form of struct request)
	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to server",
		})
	}

	//---------------------------------------------------------------------------------------------

	//sending the response

	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	r2.Decr(database.Ctx, c.IP()) //decrementing the count for rate-limiter.

	val, _ := r2.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return c.Status(fiber.StatusOK).JSON(resp)
}
