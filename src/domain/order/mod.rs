use chrono::{DateTime, Utc};
use thiserror::Error;
use crate::domain::event::{Event, OrderCreatedEvent, OrderItemSnapshot};
use crate::domain::shared::Operator;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash, Default)]
pub enum OrderStatus {
    #[default]
    Pending,
    Paid,
    Shipped,
    Delivered,
    Cancelled,
}

#[derive(Debug, Error)]
pub enum OrderError {
    #[error("order must contain at least one item")]
    EmptyItems,
    #[error("quantity must be greater than zero")]
    InvalidQuantity,
    #[error("unit price must be greater than zero")]
    InvalidUnitPrice,
    #[error("status transition is not allowed")]
    InvalidStatusTransition,
    #[error("duplicated idempotency key")]
    DuplicateIdempotencyKey,
}

pub struct OrderItemInput {
    pub product_id:   i64,
    pub product_name: String,
    pub unit_price:   i64, // cents
    pub quantity:     u32,
}

#[derive(Debug, Clone)]
pub struct OrderItem {
    pub product_id:   i64,
    pub product_name: String,
    pub unit_price:   i64,
    pub quantity:     u32,
    pub subtotal:     i64,
}
impl OrderItem {
    fn new(
        product_id: i64,
        product_name: String,
        unit_price: i64,
        quantity: u32,
    ) -> Result<Self, OrderError> {
        if quantity == 0    { return Err(OrderError::InvalidQuantity); }
        if unit_price <= 0  { return Err(OrderError::InvalidUnitPrice); }

        Ok(Self {
            product_id,
            product_name,
            unit_price,
            quantity,
            subtotal: unit_price * quantity as i64,
        })
    }
}
#[derive(Debug, Clone, Default)]
pub struct Order {
    pub id:            i64,
    pub customer_id:   i64,
    pub creator_id:    i64,
    pub creator_name:  String,
    pub modifier_id:   i64,
    pub modifier_name: String,
    pub created_at:    DateTime<Utc>,
    pub updated_at:    DateTime<Utc>,
    pub status:        OrderStatus,
    pub total_amount:  i64, // cents
    pub items:         Vec<OrderItem>,
}

impl Order {
    pub fn create(
        customer_id:   i64,
        inputs: Vec<OrderItemInput>,
        operator: Operator,
    ) -> Result<(Self, Vec<Event>), OrderError> {
        if inputs.is_empty() { return Err(OrderError::EmptyItems); }

        let items = inputs
            .into_iter()
            .map(|i| OrderItem::new(i.product_id, i.product_name, i.unit_price, i.quantity))
            .collect::<Result<Vec<OrderItem>, OrderError>>()?;
        let total_amount = items.iter().map(|i| i.subtotal).sum();
        let now = Utc::now();
        let order = Self {
            id: 0,
            customer_id,
            creator_id: operator.id,
            creator_name: operator.username,
            modifier_id: 0,
            modifier_name: String::new(),
            created_at: now,
            updated_at: now,
            status: OrderStatus::Pending,
            total_amount,
            items,
        };

        let events = vec![Event::OrderCreated(OrderCreatedEvent {
            order_id:   order.id,
            created_at: now,
            items: order.items.iter().map(|i| OrderItemSnapshot {
                product_id: i.product_id,
                quantity: i.quantity,
            }).collect(),
        })];

        Ok((order, events))
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn make_operator() -> Operator {
        Operator {
            id: 1,
            username : "alice".to_string(),
        }
    }

    fn make_inputs() -> Vec<OrderItemInput> {
        vec![
            OrderItemInput {
                product_id:   100,
                product_name: "Apple".to_string(),
                unit_price:   200,
                quantity:     3,
            },
            OrderItemInput {
                product_id:   101,
                product_name: "Banana".to_string(),
                unit_price:   500,
                quantity:     2,
            },
        ]
    }

    #[test]
    fn create_order_succeeds_with_valid_inputs() {
        let (order, events) = Order::create(42, make_inputs(), make_operator()).unwrap();
        assert_eq!(order.id, 0);
        assert_eq!(order.customer_id, 42);
        assert_eq!(order.creator_id, make_operator().id);
        assert_eq!(order.creator_name, make_operator().username);
        assert_eq!(order.modifier_id, 0);
        assert_eq!(order.modifier_name, String::new());
        assert_eq!(order.status, OrderStatus::Pending);
        assert_eq!(order.total_amount, 1600);
        assert_eq!(order.items.len(), 2);
        assert_eq!(events.len(), 1);
        assert!(matches!(events[0], Event::OrderCreated(_)))
    }

    #[test]
    fn create_order_fails_when_no_items() {
        let err = Order::create(42, vec![], make_operator()).unwrap_err();
        assert!(matches!(err, OrderError::EmptyItems));
    }

    #[test]
    fn create_order_fails_when_quantity_is_zero() {
        let bad_input = OrderItemInput {
            product_id:   100,
            product_name: "Widget".to_string(),
            unit_price:   200,
            quantity:     0,
        };
        let err = Order::create(42, vec![bad_input], make_operator()).unwrap_err();
        assert!(matches!(err, OrderError::InvalidQuantity));
    }

    #[test]
    fn create_order_fails_when_unit_price_is_zero() {
        let bad_input = OrderItemInput {
            product_id:   100,
            product_name: "Widget".to_string(),
            unit_price:   0,
            quantity:     3,
        };
        let err = Order::create(42, vec![bad_input], make_operator()).unwrap_err();
        assert!(matches!(err, OrderError::InvalidUnitPrice));
    }

}