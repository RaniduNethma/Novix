package com.novix.recommendation_service.document;

import lombok.*;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.index.Indexed;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.Map;

@Document(collation = "user_preferences")
@Getter
@Setter
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class UserPreference {

    @Id
    private String id;

    @Indexed(unique = true)
    private String userId;

    @Builder.Default
    private Map<String, Double> categoryScores = new HashMap<>();

    private LocalDateTime updatedAt;
}
