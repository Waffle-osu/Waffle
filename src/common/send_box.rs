pub struct SendBox<T> {
    pub value: T
}

unsafe impl<T> Send for SendBox<T> {

}

impl<T> SendBox<T> {
    pub unsafe fn new(value: T) -> SendBox<T> {
        SendBox { value }
    }
}