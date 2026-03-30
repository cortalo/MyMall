use thiserror::Error;

// ── Errors ────────────────────────────────────────────────────────────────────

#[derive(Debug, Error, PartialEq)]
pub enum InventoryError {
    #[error("insufficient stock for product {product_id}: requested {requested}, available {available}")]
    InsufficientStock {
        product_id: i64,
        requested:  i64,
        available:  i64,
    },
    #[error("inventory not found for product {0}")]
    NotFound(i64),
    #[error("event storage error: {0}")]
    EventStorageError(String),
}

// ── Aggregate Root ────────────────────────────────────────────────────────────

#[derive(Debug, Clone)]
pub struct Inventory {
    pub product_id: i64,
    pub quantity:   i64,
}

impl Inventory {
    pub fn new(product_id: i64, quantity: i64) -> Self {
        Self { product_id, quantity }
    }

    pub fn deduct(&mut self, quantity: i64) -> Result<(), InventoryError> {
        if self.quantity < quantity {
            return Err(InventoryError::InsufficientStock {
                product_id: self.product_id,
                requested: quantity,
                available: self.quantity,
            })
        }
        self.quantity -= quantity;
        Ok(())
    }
}

// ── Tests ─────────────────────────────────────────────────────────────────────

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn deduct_succeeds_when_stock_is_sufficient() {
        let mut inventory = Inventory::new(1, 10);
        inventory.deduct(3).unwrap();
        assert_eq!(inventory.quantity, 7);
    }

    #[test]
    fn deduct_succeeds_when_stock_is_exactly_enough() {
        let mut inventory = Inventory::new(1, 5);
        inventory.deduct(5).unwrap();
        assert_eq!(inventory.quantity, 0);
    }

    #[test]
    fn deduct_fails_when_stock_is_insufficient() {
        let mut inventory = Inventory::new(1, 3);
        let err = inventory.deduct(5).unwrap_err();
        assert_eq!(err, InventoryError::InsufficientStock {
            product_id: 1,
            requested:  5,
            available:  3,
        });
    }

    #[test]
    fn deduct_does_not_mutate_on_failure() {
        let mut inventory = Inventory::new(1, 3);
        inventory.deduct(5).unwrap_err();
        assert_eq!(inventory.quantity, 3);
    }
}