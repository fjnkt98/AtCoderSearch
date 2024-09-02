use core::fmt;
use std::error;

#[derive(Debug, Clone, Default)]
pub struct CanceledError;

impl fmt::Display for CanceledError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "canceled")
    }
}

impl error::Error for CanceledError {
    fn source(&self) -> Option<&(dyn error::Error + 'static)> {
        None
    }
}
