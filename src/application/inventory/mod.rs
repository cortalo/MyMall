use async_trait::async_trait;
use std::sync::Arc;

use crate::domain::inventory::{Inventory, InventoryError};

// ── Repository Port ───────────────────────────────────────────────────────────

#[async_trait]
pub trait InventoryRepository: Send + Sync {
    async fn save(&self, inventory: &Inventory) -> Result<(), InventoryError>;
    async fn find_by_product_id(&self, product_id: i64) -> Result<Inventory, InventoryError>;
}

// ── Service Port ──────────────────────────────────────────────────────────────

#[async_trait]
pub trait InventoryService: Send + Sync {
    async fn deduct_stock(&self, product_id: i64, quantity: i64) -> Result<(), InventoryError>;
}

// ── Service 实现 ──────────────────────────────────────────────────────────────

pub struct InventoryServiceImpl {
    repo: Arc<dyn InventoryRepository>,
}

impl InventoryServiceImpl {
    pub fn new(repo: Arc<dyn InventoryRepository>) -> Self {
        Self { repo }
    }
}

#[async_trait]
impl InventoryService for InventoryServiceImpl {
    async fn deduct_stock(&self, product_id: i64, quantity: i64) -> Result<(), InventoryError> {
        let mut inventory = self.repo.find_by_product_id(product_id).await?;
        inventory.deduct(quantity)?;
        self.repo.save(&inventory).await?;
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

    // ── Mock ──────────────────────────────────────────────────────────────────

    struct MockInventoryRepository {
        save_fn: Box<dyn Fn(&Inventory) -> Result<(), InventoryError> + Send + Sync>,
        find_by_product_id_fn: Box<dyn Fn(i64) -> Result<Inventory, InventoryError> + Send + Sync>,
    }

    impl MockInventoryRepository {
        fn new() -> Self {
            Self {
                find_by_product_id_fn: Box::new(|_| panic!("find_by_product_id should not be called")),
                save_fn: Box::new(|_| panic!("save should not be called")),
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
    }

    #[async_trait]
    impl InventoryRepository for MockInventoryRepository {
        async fn save(&self, inventory: &Inventory) -> Result<(), InventoryError> {
            (self.save_fn)(inventory)
        }
        async fn find_by_product_id(&self, product_id: i64) -> Result<Inventory, InventoryError> {
            (self.find_by_product_id_fn)(product_id)
        }

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

        let service = InventoryServiceImpl::new(Arc::new(repo));
        service.deduct_stock(1, 3).await?;

        Ok(())
    }

    #[tokio::test]
    async fn deduct_stock_fails_when_product_not_found() {
        let repo = MockInventoryRepository::new()
            .with_find_by_product_id(|product_id| Err(InventoryError::NotFound(product_id)));

        let service = InventoryServiceImpl::new(Arc::new(repo));
        let err = service.deduct_stock(1, 3).await.unwrap_err();

        assert_eq!(err, InventoryError::NotFound(1));
    }

    #[tokio::test]
    async fn deduct_stock_fails_when_stock_is_insufficient() {
        let repo = MockInventoryRepository::new()
            .with_find_by_product_id(|_| Ok(Inventory::new(1, 2)))
            .with_save(|_| panic!("save should not be called when deduct fails"));

        let service = InventoryServiceImpl::new(Arc::new(repo));
        let err = service.deduct_stock(1, 5).await.unwrap_err();

        assert_eq!(err, InventoryError::InsufficientStock {
            product_id: 1,
            requested:  5,
            available:  2,
        });
    }
}