package com.novix.notification_service.repository;

import com.novix.notification_service.entity.Notification;
import com.novix.notification_service.enums.NotificationType;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

@Repository
public interface NotificationRepository extends JpaRepository<Notification, Long> {
    Page<Notification> findByUserId(String userId, Pageable pageable);

    Page<Notification> findByUserIdAndIsRead(String userId, Boolean isRead, Pageable pageable);

    Page<Notification> findByUserIdAndType(String userId, NotificationType type, Pageable pageable);

    Long countByUserIdAndIsRead (String userId, Boolean isRead);

    @Modifying
    @Query("UPDATE Notification n SET n.isRead = true, n.status = 'READ' WHERE n.userId = :userId AND n.isRead = false")
    void markAllAsReadByUserId(String userId);
}
