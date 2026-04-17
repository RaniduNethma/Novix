package com.novix.user_service.dto.response;

import java.time.LocalDateTime;
import java.util.Set;
import com.fasterxml.jackson.annotation.JsonProperty;
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
	
	@JsonProperty("isActive")
	private boolean active;
	
	@JsonProperty("isEmailVerified")
	private boolean emailVerified;
	
	private Set<String> roles;
	private LocalDateTime createdAt;
}
