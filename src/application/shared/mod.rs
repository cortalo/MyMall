use async_trait::async_trait;
use crate::domain::event::Event;

#[async_trait]
pub trait EventBus: Send + Sync {
    async fn publish(&self, event: Event) -> Result<(), EventBusError>;
}

#[derive(Debug, thiserror::Error)]
pub enum EventBusError {
    #[error("failed to publish event: {0}")]
    PublishFailed(String),
}

// ── ProcessedEventRepository ──────────────────────────────────────────────────

#[async_trait]
pub trait ProcessedEventRepository: Send + Sync {
    async fn is_processed(&self, event_type: &str, event_id: i64) -> Result<bool, ProcessedEventError>;
    async fn mark_as_processed(&self, event_type: &str, event_id: i64) -> Result<(), ProcessedEventError>;
}

#[derive(Debug, thiserror::Error)]
pub enum ProcessedEventError {
    #[error("failed to check processed event: {0}")]
    StorageFailed(String),
}