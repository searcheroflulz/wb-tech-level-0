package generator

import (
	"fmt"
	"math/rand"
	"time"
	"wb-tech-level-0/internal/model"
)

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) GenerateOrder() *model.Order {
	// Генерация случайных данных
	orderUID := fmt.Sprintf("b563feb7b2b84b6test%d", rand.Intn(1000))
	trackNumber := fmt.Sprintf("WBILMTESTTRACK%d", rand.Intn(1000))
	// Генерация случайного числа элементов в слайсе Items
	numOrderItems := rand.Intn(5) + 1 // Минимум 1 элемент
	orderItems := make([]model.Items, numOrderItems)

	for i := 0; i < numOrderItems; i++ {
		orderItems[i] = model.Items{
			ChrtID:      rand.Intn(1000000),
			TrackNumber: trackNumber,
			Price:       rand.Intn(500),
			RID:         fmt.Sprintf("rid%d", rand.Intn(1000)),
			Name:        "Mascaras",
			Sale:        rand.Intn(50),
			Size:        fmt.Sprintf("%d", rand.Intn(45)),
			TotalPrice:  rand.Intn(500),
			NMID:        rand.Intn(1000000),
			Brand:       "Vivienne Sabo",
			Status:      rand.Intn(300),
		}
	}

	return &model.Order{
		OrderUID:          orderUID,
		TrackNumber:       trackNumber,
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          fmt.Sprintf("%d", rand.Intn(9)),
		SmID:              rand.Intn(99),
		DateCreated:       time.Now(),
		OofShard:          fmt.Sprintf("%d", rand.Intn(5)),
		Delivery: model.Delivery{
			Name:    "Test Testov",
			Phone:   fmt.Sprintf("+972%09d", rand.Intn(1000000000)),
			Zip:     fmt.Sprintf("%07d", rand.Intn(10000000)),
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: model.Payment{
			Transaction:  orderUID,
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       rand.Intn(2000),
			PaymentDate:  int(time.Now().Unix()),
			Bank:         "alpha",
			DeliveryCost: rand.Intn(1000),
			GoodsTotal:   rand.Intn(500),
			CustomFee:    0,
		},
		Items: orderItems,
	}
}
