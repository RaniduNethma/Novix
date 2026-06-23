package com.novix.notification_service.service.impl;

import com.novix.notification_service.service.EmailService;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import lombok.extern.slf4j.Slf4j;
import org.springframework.mail.SimpleMailMessage;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
@Slf4j
public class EmailServiceImpl implements EmailService {
    private final JavaMailSender mailSender;

    @Value("${notification.email.from}")
    private String fromEmail;

    @Value("${notification.email.enabled}")
    private boolean emailEnabled;

    @Override
    public void sendEmail(String to, String subject, String body) {
        if (!emailEnabled){
            log.info("Email disabled, skipping send to: {}", to);
            return;
        }

        try{
            SimpleMailMessage message = new SimpleMailMessage();
            message.setFrom(fromEmail);
            message.setTo(to);
            message.setSubject(subject);
            message.setText(body);
            mailSender.send(message);
            log.info("Email sent to: {}", to);
        } catch (Exception e){
            log.error("Failed to send email to: {} - {}", to, e.getMessage());
            throw new RuntimeException("Email sending failed: " + e.getMessage());
        }
    }
}
