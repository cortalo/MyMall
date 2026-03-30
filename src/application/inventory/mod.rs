use async_trait::async_trait;
use std::sync::Arc;
use crate::application::shared::{ProcessedEventError, ProcessedEventRepository};
use crate::domain::event::OrderCreatedEvent;
use crate::domain::inventory::{Inventory, InventoryError};

const ORDER_CREATED_EVENT_TYPE: &str = "order.created";

// ── Repository Port ───────────────────────────────────────────────────────────

#[async_trait]
pub trait InventoryRepository: Send + Sync {
    async fn save(&self, inventory: &Inventory) -> Result<(), InventoryError>;
    async fn find_by_product_id(&self, product_id: i64) -> Result<Inventory, InventoryError>;
    async fn deduct_stock(&self, product_id: i64, quantity: u32) -> Result<(), InventoryError>;
    async fn deduct_stocks_batch(&self, deductions: &[StockDeduction]) -> Result<(), InventoryError>;
}
pub struct StockDeduction {
    pub product_id: i64,
    pub quantity:   u32,
}

// ── Service Port ──────────────────────────────────────────────────────────────

#[async_trait]
pub trait InventoryService: Send + Sync {
    async fn deduct_stock(&self, product_id: i64, quantity: u32) -> Result<(), InventoryError>;
    async fn handle_order_created(&self, event: OrderCreatedEvent) -> Result<(), InventoryError>;
}

// ── Service 实现 ──────────────────────────────────────────────────────────────

pub struct InventoryServiceImpl {
    repo: Arc<dyn InventoryRepository>,
    processed_event_repo:   Arc<dyn ProcessedEventRepository>,
}

impl InventoryServiceImpl {
    pub fn new(
        repo: Arc<dyn InventoryRepository>,
        processed_event_repo: Arc<dyn ProcessedEventRepository>,
    ) -> Self {
        Self { repo, processed_event_repo }
    }
}

#[async_trait]
impl InventoryService for InventoryServiceImpl {
    async fn deduct_stock(&self, product_id: i64, quantity: u32) -> Result<(), InventoryError> {
        let mut inventory = self.repo.find_by_product_id(product_id).await?;
        inventory.deduct(quantity as i64)?;
        self.repo.save(&inventory).await?;
        Ok(())
    }

    async fn handle_order_created(&self, event: OrderCreatedEvent) -> Result<(), InventoryError> {
        let already_processed = self.processed_event_repo
            .is_processed(ORDER_CREATED_EVENT_TYPE, event.order_id)
            .await
            .map_err(|e| InventoryError::EventStorageError(e.to_string()))?;
        if already_processed {
            return Ok(());
        }

        // TODO: how to support transaction
        for item in event.items {
            self.deduct_stock(item.product_id, item.quantity).await?;
        }

        self.processed_event_repo
            .mark_as_processed(ORDER_CREATED_EVENT_TYPE, event.order_id)
            .await
            .map_err(|e| InventoryError::EventStorageError(e.to_string()))?;

        Ok(())
    }
}

// ── Tests ─────────────────────────────────────────────────────────────────────

#[cfg(test)]
mod tests {
    use super::*;
    use crate::domain::inventory::{Inventory, InventoryError};
    use async_trait::async_trait;
    use std::sync::Arc;
    use chrono::Utc;
    use crate::domain::event::OrderItemSnapshot;
    // ── Mock ──────────────────────────────────────────────────────────────────

    struct MockInventoryRepository {
        save_fn: Box<dyn Fn(&Inventory) -> Result<(), InventoryError> + Send + Sync>,
        find_by_product_id_fn: Box<dyn Fn(i64) -> Result<Inventory, InventoryError> + Send + Sync>,
        deduct_stock_fn: Box<dyn Fn(i64, u32) -> Result<(), InventoryError> + Send + Sync>,
        deduct_stocks_batch_fn: Box<dyn Fn(&[StockDeduction]) -> Result<(), InventoryError> + Send + Sync>,
    }

    impl MockInventoryRepository {
        fn new() -> Self {
            Self {
                find_by_product_id_fn: Box::new(|_| panic!("find_by_product_id should not be called")),
                save_fn: Box::new(|_| panic!("save should not be called")),
                deduct_stock_fn: Box::new(|_, _| panic!("deduct stock should not be called")),
                deduct_stocks_batch_fn: Box::new(|_| panic!("deduct stocks batch should not be called")),
            }
        }

        fn with_find_by_product_id(mut self, f: impl Fn(i64) -> Result<Inventory, InventoryError> + Send + Sync + 'static) -> Self {
            self.find_by_product_id_fn = Box::new(f);
            self
        }
        fn with_save(mut self, f: impl Fn(&Inventory) -> Result<(), InventoryError> + Send + Sync + 'static) -> Self {
            self.save_fn = Box::new(f);
            self
        }
        fn with_deduct_stock(mut self, f: impl Fn(i64, u32) -> Result<(), InventoryError> + Send + Sync + 'static) -> Self {
            self.deduct_stock_fn = Box::new(f);
            self
        }
        fn with_deduct_stocks_batch(mut self, f: impl Fn(&[StockDeduction]) -> Result<(), InventoryError> + Send + Sync + 'static) -> Self {
            self.deduct_stocks_batch_fn = Box::new(f);
            self
        }

    }

    #[async_trait]
    impl InventoryRepository for MockInventoryRepository {
        async fn save(&self, inventory: &Inventory) -> Result<(), InventoryError> {
            (self.save_fn)(inventory)
        }
        async fn find_by_product_id(&self, product_id: i64) -> Result<Inventory, InventoryError> {
            (self.find_by_product_id_fn)(product_id)
        }
        async fn deduct_stock(&self, product_id: i64, quantity: u32) -> Result<(), InventoryError> {
            (self.deduct_stock_fn)(product_id, quantity)
        }

        async fn deduct_stocks_batch(&self, deductions: &[StockDeduction]) -> Result<(), InventoryError> {
            (self.deduct_stocks_batch_fn)(deductions)
        }
    }

    // ── Mock ProcessedEventRepository ─────────────────────────────────────────

    struct MockProcessedEventRepository {
        is_processed_fn:    Box<dyn Fn(&str, i64) -> Result<bool, ProcessedEventError> + Send + Sync>,
        mark_processed_fn:  Box<dyn Fn(&str, i64) -> Result<(), ProcessedEventError> + Send + Sync>,
    }

    impl MockProcessedEventRepository {
        fn new() -> Self {
            Self {
                is_processed_fn:   Box::new(|_, _| Ok(false)),
                mark_processed_fn: Box::new(|_, _| Ok(())),
            }
        }

        fn with_is_processed(mut self, f: impl Fn(&str, i64) -> Result<bool, ProcessedEventError> + Send + Sync + 'static) -> Self {
            self.is_processed_fn = Box::new(f);
            self
        }
        fn with_mark_processed(mut self, f: impl Fn(&str, i64) -> Result<(), ProcessedEventError> + Send + Sync + 'static) -> Self {
            self.mark_processed_fn = Box::new(f);
            self
        }
    }

    #[async_trait]
    impl ProcessedEventRepository for MockProcessedEventRepository {
        async fn is_processed(&self, event_type: &str, event_id: i64) -> Result<bool, ProcessedEventError> {
            (self.is_processed_fn)(event_type, event_id)
        }
        async fn mark_as_processed(&self, event_type: &str, event_id: i64) -> Result<(), ProcessedEventError> {
            (self.mark_processed_fn)(event_type, event_id)
        }
    }

    fn make_order_created_event(order_id: i64, items: Vec<(i64, u32)>) -> OrderCreatedEvent {
        OrderCreatedEvent {
            order_id,
            created_at: Utc::now(),
            items:      items.into_iter()
                .map(|(product_id, quantity)| OrderItemSnapshot { product_id, quantity })
                .collect(),
        }
    }

    fn make_service(
        repo: MockInventoryRepository,
        processed_repo: MockProcessedEventRepository,
    ) -> InventoryServiceImpl {
        InventoryServiceImpl::new(Arc::new(repo), Arc::new(processed_repo))
    }


    // ── Tests ─────────────────────────────────────────────────────────────────

    #[tokio::test]
    async fn deduct_stock_succeeds() -> Result<(), InventoryError> {
        let repo = MockInventoryRepository::new()
            .with_find_by_product_id(|product_id| {
                assert_eq!(product_id, 1);
                Ok(Inventory::new(1, 10))
            })
            .with_save(|inventory| {
                assert_eq!(inventory.product_id, 1);
                assert_eq!(inventory.quantity, 7); // 10 - 3
                Ok(())
            });

        let service = make_service(repo, MockProcessedEventRepository::new());
        service.deduct_stock(1, 3).await?;

        Ok(())
    }

    // #[tokio::test]
    // async fn deduct_stock_fails_when_product_not_found() {
    //     let repo = MockInventoryRepository::new()
    //         .with_find_by_product_id(|product_id| Err(InventoryError::NotFound(product_id)));
    //
    //     let service = make_service(repo, MockProcessedEventRepository::new());
    //     let err = service.deduct_stock(1, 3).await.unwrap_err();
    //
    //     assert_eq!(err, InventoryError::NotFound(1));
    // }
    //
    // #[tokio::test]
    // async fn deduct_stock_fails_when_stock_is_insufficient() {
    //     let repo = MockInventoryRepository::new()
    //         .with_find_by_product_id(|_| Ok(Inventory::new(1, 2)))
    //         .with_save(|_| panic!("save should not be called when deduct fails"));
    //
    //     let service = make_service(repo, MockProcessedEventRepository::new());
    //     let err = service.deduct_stock(1, 5).await.unwrap_err();
    //
    //     assert_eq!(err, InventoryError::InsufficientStock {
    //         product_id: 1,
    //         requested:  5,
    //         available:  2,
    //     });
    // }

    #[tokio::test]
    async fn handle_order_created_deducts_stock_for_all_items() {
        let calls = Arc::new(std::sync::Mutex::new(vec![]));
        let calls_clone = calls.clone();

        let repo = MockInventoryRepository::new()
            .with_find_by_product_id(|product_id| Ok(Inventory::new(product_id, 100)))
            .with_save(move |inventory| {
                calls_clone.lock().unwrap().push((inventory.product_id, inventory.quantity));
                Ok(())
            });

        let service = make_service(repo, MockProcessedEventRepository::new());
        service.handle_order_created(make_order_created_event(1, vec![(1, 3), (2, 5)])).await.unwrap();

        let calls = calls.lock().unwrap();
        assert_eq!(calls.len(), 2);
        assert_eq!(calls[0], (1, 97)); // 100 - 3
        assert_eq!(calls[1], (2, 95)); // 100 - 5
    }

    // #[tokio::test]
    // async fn handle_order_created_stops_on_first_failure() {
    //     let call_count = Arc::new(std::sync::Mutex::new(0));
    //     let call_count_clone = call_count.clone();
    //
    //     let repo = MockInventoryRepository::new()
    //         .with_find_by_product_id(move |product_id| {
    //             *call_count_clone.lock().unwrap() += 1;
    //             if product_id == 2 {
    //                 Ok(Inventory::new(2, 0)) // 库存不足
    //             } else {
    //                 Ok(Inventory::new(product_id, 100))
    //             }
    //         })
    //         .with_save(|_| Ok(()));
    //
    //     let service = make_service(repo, MockProcessedEventRepository::new());
    //     let err = service
    //         .handle_order_created(make_order_created_event(1, vec![(1, 3), (2, 5), (3, 1)]))
    //         .await
    //         .unwrap_err();
    //
    //     assert_eq!(*call_count.lock().unwrap(), 2); // product 3 没有被处理
    //     assert!(matches!(err, InventoryError::InsufficientStock { product_id: 2, .. }));
    // }

    #[tokio::test]
    async fn handle_order_created_skips_if_already_processed() {
        let repo = MockInventoryRepository::new()
            .with_find_by_product_id(|_| panic!("should not query inventory for duplicate event"))
            .with_save(|_| panic!("should not save inventory for duplicate event"));

        let processed_repo = MockProcessedEventRepository::new()
            .with_is_processed(|_, _| Ok(true)); // 已经处理过

        let service = make_service(repo, processed_repo);
        service.handle_order_created(make_order_created_event(1, vec![(1, 3)])).await.unwrap();
    }

    // #[tokio::test]
    // async fn handle_order_created_marks_event_as_processed() {
    //     let marked = Arc::new(std::sync::Mutex::new(vec![]));
    //     let marked_clone = marked.clone();
    //
    //     let repo = MockInventoryRepository::new()
    //         .with_deduct_stocks_batch(|_| Ok(()));
    //
    //     let processed_repo = MockProcessedEventRepository::new()
    //         .with_mark_processed(move |event_type, event_id| {
    //             marked_clone.lock().unwrap().push((event_type.to_string(), event_id));
    //             Ok(())
    //         });
    //
    //     let service = make_service(repo, processed_repo);
    //     service.handle_order_created(make_order_created_event(42, vec![(1, 3)])).await.unwrap();
    //
    //     let marked = marked.lock().unwrap();
    //     assert_eq!(marked.len(), 1);
    //     assert_eq!(marked[0], (ORDER_CREATED_EVENT_TYPE.to_string(), 42));
    // }
}