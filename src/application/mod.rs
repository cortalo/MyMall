// TODO
// check application read/write race conditions
// group multiple read/write to single interface so that enable transaction in infra
// publish create order event with order_id
mod order;
pub mod inventory;
pub mod shared;