package com.novix.user_service.dto.response;

import java.time.LocalDateTime;
import java.util.Set;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class UserResponse {
	private Long id;
	private String firstName;
	private String lastName;
	private String userName;
	private String email;
	private String profilePicture;
	private boolean isActive;
	private boolean isEmailVerified;
	private Set<String> roles;
	private LocalDateTime createdAt;
}
