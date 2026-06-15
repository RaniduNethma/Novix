package com.novix.notification_service.exception;

public class NotificationNotFoundException extends RuntimeException{
    public NotificationNotFoundException(String message){
        super(message);
    }
}
