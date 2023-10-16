package routes

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type CoinDto struct {
	CoinId        uint    `json:"coin_id"`
	Id            string  `json:"id"`
	Large         string  `json:"large"`
	MarketCapRank uint    `json:"market_cap_rank"`
	Name          string  `json:"name"`
	PriceBtc      float64 `json:"price_btc"`
	Score         uint    `json:"score"`
	Slug          string  `json:"slug"`
	Small         string  `json:"small"`
	Symbol        string  `json:"symbol"`
	Thumb         string  `json:"thumb"`
	Abc           string  `json:"abc"`
}

type CoinItemDto struct {
	Item CoinDto `json:"item"`
}

type TrendingCoinsDto struct {
	Coins []CoinItemDto `json:"coins"`
}

// GetTrendingCoins godoc
// @Summary      List of trending coins
// @Description  get list of trending coin
// @Tags         coins
// @Accept       json
// @Produce      json
// @Success      200  {object}   TrendingCoinsDto
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /api/coin [get]
func GetTrendingCoins(c *fiber.Ctx) error {
	// declare agent
	agent := fiber.Get("https://api.coingecko.com/api/v3/search/trending")
	// get response and status
	statusCode, body, errs := agent.Bytes()
	// if error, return error
	if len(errs) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}
	var responseData TrendingCoinsDto
	err := json.Unmarshal(body, &responseData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": err,
		})
	}
	fmt.Println(responseData)
	return c.Status(statusCode).JSON(responseData)
}
