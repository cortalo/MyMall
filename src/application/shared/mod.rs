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