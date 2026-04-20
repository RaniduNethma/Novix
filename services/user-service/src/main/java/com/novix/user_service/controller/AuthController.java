package com.novix.user_service.controller;

import com.novix.user_service.dto.request.LoginRequest;
import com.novix.user_service.dto.request.RegisterRequest;
import com.novix.user_service.dto.response.AuthResponse;
import com.novix.user_service.service.AuthService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1/auth")
@RequiredArgsConstructor
public class AuthController {
	private final AuthService authService;
	
	@PostMapping("/register")
	public ResponseEntity<AuthResponse> register(
			@Valid @RequestBody RegisterRequest request){
		return ResponseEntity
				.status(HttpStatus.CREATED)
				.body(authService.register(request));
	}
	
	@PostMapping("/login")
	public ResponseEntity<AuthResponse> login(
			@Valid @RequestBody LoginRequest request){
		return ResponseEntity.ok(authService.login(request));
	}
	
	@PostMapping("/refresh-token")
	public ResponseEntity<AuthResponse> refreshToken(
			@RequestHeader("Refresh-Token") String refreshToken){
		return ResponseEntity.ok(authService.refreshToken(refreshToken));
	}
	
	@PostMapping("/logout")
	public ResponseEntity<String> logout(
			@RequestHeader("Authorization") String authHeader){
		String token = authHeader.substring(7);
		return ResponseEntity.ok("Logged out successfully");
	}
}
