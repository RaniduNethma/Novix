package com.novix.notification_service.mapper;

import com.novix.notification_service.dto.rsponse.NotificationResponse;
import com.novix.notification_service.entity.Notification;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface NotificationMapper {

    @Mapping(target = "type", expression = "java(notification.getType().name())")
    @Mapping(target = "status", expression = "java(notification.getStatus().name())")
    NotificationResponse toResponse(Notification notification);
}
