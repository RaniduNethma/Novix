package com.novix.user_service.service;

import com.novix.user_service.dto.request.LoginRequest;
import com.novix.user_service.dto.request.RegisterRequest;
import com.novix.user_service.dto.response.AuthResponse;

public interface AuthService {
	AuthResponse register(RegisterRequest request);
	AuthResponse login(LoginRequest request);
	AuthResponse refreshToken(String refreshToken);
	void logout(String email);
}
