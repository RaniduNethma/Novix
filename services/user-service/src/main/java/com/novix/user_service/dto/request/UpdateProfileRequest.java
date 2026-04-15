package com.novix.user_service.dto.request;

import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class UpdateProfileRequest {
	
	@Size(min = 1, message = "First name cannot be empty")
	private String firstName;
	
	@Size(min = 1, message = "Last name cannot be empty")
	private String lastName;
	
	private String profilePicture;
}
