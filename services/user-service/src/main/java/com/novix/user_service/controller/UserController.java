package com.novix.user_service.controller;

import com.novix.user_service.dto.request.UpdateProfileRequest;
import com.novix.user_service.dto.response.UserResponse;
import com.novix.user_service.service.UserService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1/users")
@RequiredArgsConstructor
public class UserController {
	private final UserService userService;
	
	@GetMapping("/me")
	public ResponseEntity<UserResponse> getCurrentUser(
			@AuthenticationPrincipal UserDetails userDetails){
		return ResponseEntity.ok(
				userService.getUserByEmail(userDetails.getUsername())
		);
	}
	
	@GetMapping("/{id}")
	@PreAuthorize("hasRole('ADMIN')")
	public ResponseEntity<UserResponse> getUserById(@PathVariable Long id) {
		return ResponseEntity.ok(userService.getUserById(id));
	}
	
	@PutMapping("/{id}")
	public ResponseEntity<UserResponse> updateProfile(
			@PathVariable Long id,
			@Valid @RequestBody UpdateProfileRequest request) {
		return ResponseEntity.ok(userService.updateProfile(id, request));
	}
	
	@DeleteMapping("/{id}")
	@PreAuthorize("hasRole('ADMIN')")
	public ResponseEntity<String> deleteUser(@PathVariable Long id) {
		userService.deleteUser(id);
		return ResponseEntity.ok("User deleted successfully");
	}
}
