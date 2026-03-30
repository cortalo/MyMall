use chrono::{DateTime, Utc};

#[derive(Debug, Clone)]
pub enum Event {
    OrderCreated(OrderCreatedEvent),
}

#[derive(Debug, Clone)]
pub struct OrderCreatedEvent {
    pub order_id: i64,
    pub created_at: DateTime<Utc>,
}