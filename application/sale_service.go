package application

import (
	"MyMall/domain"
	"MyMall/domain/shared"
	"context"
	"errors"
	"time"

	"github.com/samber/lo"
)

type SaleService interface {
	CreateOrder(
		ctx context.Context,
		customerID int64,
		items []domain.SaleOrderItem,
		operator shared.Operator,
		idempotencyKey string,
	) (*domain.Order, error)
}

type SaleRepository interface {
	FindByID(ctx context.Context, id int64) (*domain.Sale, error)
}

type OrderReadRepository interface {
	// CountByUserAndProduct implementation needs FOR UPDATE to avoid write skew
	CountByUserAndProduct(ctx context.Context, userID int64, productID int64) (int, error)
	FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Order, error)
}

type OrderCreator interface {
	CreateOrder(ctx context.Context, customerID int64, inputs []domain.OrderItemInput, operator shared.Operator, idempotencyKey string) (*domain.Order, error)
}

type saleService struct {
	saleRepo      SaleRepository
	orderReadRepo OrderReadRepository
	orderCreator  OrderCreator
}

func NewSaleService(
	saleRepo SaleRepository,
	orderReadRepo OrderReadRepository,
	orderCreator OrderCreator,
) SaleService {
	return &saleService{
		saleRepo:      saleRepo,
		orderReadRepo: orderReadRepo,
		orderCreator:  orderCreator,
	}
}

func (s *saleService) CreateOrder(
	ctx context.Context,
	customerID int64,
	items []domain.SaleOrderItem,
	operator shared.Operator,
	idempotencyKey string,
) (*domain.Order, error) {
	order, err := s.orderReadRepo.FindByIdempotencyKey(ctx, idempotencyKey)
	if err != nil && !errors.Is(err, domain.ErrOrderNotFound) {
		return nil, err
	}
	if order != nil {
		return order, nil
	}
	orderItemInputs, err := lo.MapErr(items, func(item domain.SaleOrderItem, _ int) (domain.OrderItemInput, error) {
		sale, err := s.saleRepo.FindByID(ctx, item.SaleID)
		if err != nil {
			return domain.OrderItemInput{}, err
		}
		if !sale.IsActive(time.Now()) {
			return domain.OrderItemInput{}, domain.ErrSaleNotActive
		}
		currentCount, err := s.orderReadRepo.CountByUserAndProduct(ctx, customerID, sale.ProductID)
		if err != nil {
			return domain.OrderItemInput{}, err
		}
		if err := sale.CheckPurchaseLimit(currentCount, item.Quantity); err != nil {
			return domain.OrderItemInput{}, err
		}
		return domain.OrderItemInput{
			ProductID:   sale.ProductID,
			ProductName: sale.ProductName,
			UnitPrice:   sale.UnitPrice,
			Quantity:    item.Quantity,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return s.orderCreator.CreateOrder(ctx, customerID, orderItemInputs, operator, idempotencyKey)
}
