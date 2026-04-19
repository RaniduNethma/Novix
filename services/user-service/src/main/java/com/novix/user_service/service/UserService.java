package com.novix.user_service.service;

import com.novix.user_service.dto.request.UpdateProfileRequest;
import com.novix.user_service.dto.response.UserResponse;

public interface UserService {
	UserResponse getUserById(Long id);
	UserResponse getUserByEmail(String email);
	UserResponse updateProfile(Long id, UpdateProfileRequest request);
	void deleteUser(Long id);
}
