use async_trait::async_trait;
use std::sync::Arc;

use crate::domain::order::{Order, OrderError, OrderItemInput};
use crate::domain::shared::Operator;

#[async_trait]
pub trait OrderRepository: Send + Sync {
    async fn save(&self, order: &Order, idempotency_key: &str) -> Result<(), OrderError>;
    async fn find_by_id(&self, id: i64) -> Result<Order, OrderError>;
    async fn find_by_idempotency_key(&self, idempotency_key: &str) -> Result<Order, OrderError>;
}

#[async_trait]
pub trait OrderService: Send + Sync {
    async fn create_order(
        &self,
        customer_id: i64,
        inputs: Vec<OrderItemInput>,
        operator: Operator,
        idempotency_key: &str,
    ) -> Result<Order, OrderError>;
}

pub struct OrderServiceImpl {
    repo: Arc<dyn OrderRepository>,
}

impl OrderServiceImpl {
    pub fn new(repo: Arc<dyn OrderRepository>) -> Self {
        Self { repo }
    }
}

#[async_trait]
impl OrderService for OrderServiceImpl {
    async fn create_order(
        &self,
        customer_id: i64,
        inputs: Vec<OrderItemInput>,
        operator: Operator,
        idempotency_key: &str
    ) -> Result<Order, OrderError> {
        let (order, _events) = Order::create(customer_id, inputs, operator)?;
        match self.repo.save(&order, idempotency_key).await {
            Ok(_) => Ok(order),
            Err(OrderError::DuplicateIdempotencyKey) => {
                self.repo.find_by_idempotency_key(idempotency_key).await
            }
            Err(e) => Err(e),
        }
    }
}

#[cfg(test)]
mod tests {
    use crate::domain::order::OrderStatus;
use super::*;

    struct MockOrderRepository {
        save_fn: Box<dyn Fn(&Order, &str) -> Result<(), OrderError> + Send + Sync>,
        find_by_id_fn: Box<dyn Fn(i64) -> Result<Order, OrderError> + Send + Sync>,
        find_by_idempotency_key_fn: Box<dyn Fn(&str) -> Result<Order, OrderError> + Send + Sync>,
    }

    impl MockOrderRepository {
        fn new() -> Self {
            Self {
                save_fn: Box::new(|_, _| Ok(())),
                find_by_idempotency_key_fn: Box::new(|_| panic!("find_by_idempotency_key should not be called")),
                find_by_id_fn: Box::new(|_| panic!("find_by_id should not be called")),
            }
        }

        fn with_save(mut self, f: impl Fn(&Order, &str) -> Result<(), OrderError> + Send + Sync + 'static) -> Self {
            self.save_fn = Box::new(f);
            self
        }

        fn with_find_by_idempotency_key(mut self, f: impl Fn(&str) -> Result<Order, OrderError> + Send + Sync + 'static) -> Self {
            self.find_by_idempotency_key_fn = Box::new(f);
            self
        }
    }

    #[async_trait]
    impl OrderRepository for MockOrderRepository {
        async fn save(&self, order: &Order, idempotency_key: &str) -> Result<(), OrderError> {
            (self.save_fn)(order, idempotency_key)
        }
        async fn find_by_id(&self, id: i64) -> Result<Order, OrderError> {
            (self.find_by_id_fn)(id)
        }
        async fn find_by_idempotency_key(&self, idempotency_key: &str) -> Result<Order, OrderError> {
            (self.find_by_idempotency_key_fn)(idempotency_key)
        }
    }

    fn make_operator() -> Operator {
        Operator { id: 1, username: "test_user".to_string() }
    }
    fn make_inputs() -> Vec<OrderItemInput> {
        vec![
            OrderItemInput { product_id: 1, product_name: "Apple".to_string(), unit_price: 500, quantity: 2 },
        ]
    }

    #[tokio::test]
    async fn create_order_succeeds() -> Result<(), OrderError> {
        let repo = MockOrderRepository::new()
            .with_save(|_, idempotency_key| {
                assert_eq!(idempotency_key, "random-key-72");
                Ok(())
            });

        let service = OrderServiceImpl::new(Arc::new(repo));

        let order = service
            .create_order(42, make_inputs(), make_operator(), "random-key-72")
            .await?;

        assert_eq!(order.customer_id, 42);
        assert_eq!(order.status, OrderStatus::Pending);
        assert_eq!(order.total_amount, 1000);

        Ok(())
    }

    #[tokio::test]
    async fn create_order_idempotent_retry() -> Result<(), OrderError> {
        let existing_order = Order {
            id: 99,
            customer_id: 42,
            status: OrderStatus::Pending,
            total_amount: 1000,
            ..Order::default()
        };

        let repo = MockOrderRepository::new()
            .with_save(|_, _| Err(OrderError::DuplicateIdempotencyKey))
            .with_find_by_idempotency_key(move |idempotency_key| {
                assert_eq!(idempotency_key, "idempotency-key-1");
                Ok(existing_order.clone())
            });

        let service = OrderServiceImpl::new(Arc::new(repo));

        let order = service
            .create_order(42, make_inputs(), make_operator(), "idempotency-key-1")
            .await?;

        assert_eq!(order.id, 99);
        assert_eq!(order.customer_id, 42);

        Ok(())
    }
}