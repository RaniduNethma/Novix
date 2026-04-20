package com.novix.user_service.service.impl;

import java.util.HashSet;
import java.util.Set;
import com.novix.user_service.dto.request.LoginRequest;
import com.novix.user_service.dto.request.RegisterRequest;
import com.novix.user_service.dto.response.AuthResponse;
import com.novix.user_service.entity.Role;
import com.novix.user_service.entity.User;
import com.novix.user_service.enums.RoleType;
import com.novix.user_service.exception.EmailAlreadyExistsException;
import com.novix.user_service.exception.UserNotFoundException;
import com.novix.user_service.mapper.UserMapper;
import com.novix.user_service.repository.RoleRepository;
import com.novix.user_service.repository.UserRepository;
import com.novix.user_service.security.JwtService;
import com.novix.user_service.service.AuthService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@RequiredArgsConstructor
@Slf4j
public class AuthServiceImpl implements AuthService {
	private final UserRepository userRepository;
	private final RoleRepository roleRepository;
	private final PasswordEncoder passwordEncoder;
	private final JwtService jwtService;
	private final AuthenticationManager authenticationManager;
	private final UserDetailsService userDetailsService;
	private final UserMapper userMapper;
	
	@Override
	@Transactional
	public AuthResponse register(RegisterRequest request) {
		
		// Check email already exists
		if (userRepository.existsByEmail(request.getEmail())){
			throw new EmailAlreadyExistsException(
					"Email already registered: " + request.getEmail()
			);
		}
		
		// Check email already exists
		if (userRepository.existsByUserName(request.getUserName())){
			throw new EmailAlreadyExistsException(
					"Username already taken: " + request.getUserName()
			);
		}
		
		// Assign default role "USER"
		Role userRole = roleRepository.findByName(RoleType.ROLE_USER).orElseThrow(() -> new RuntimeException(
				"Default role not found. Please seed roles first."
		));
		Set<Role> roles = new HashSet<>();
		roles.add(userRole);
		
		// Build and save user
		User user = User.builder()
				.firstName(request.getFirstName())
				.lastName(request.getLastName())
				.userName(request.getUserName())
				.email(request.getEmail())
				.password(passwordEncoder.encode(request.getPassword()))
				.roles(roles)
				.isActive(true)
				.isEmailVerified(false)
				.build();
		User savedUser = userRepository.save(user);
		
		// Generate tokens
		UserDetails userDetails = userDetailsService.loadUserByUsername(savedUser.getEmail());
		String accessToken = jwtService.generateAccessToken(userDetails);
		String refreshToken = jwtService.generateRefreshToken(userDetails);
		
		// Save refresh token
		savedUser.setRefreshToken(refreshToken);
		userRepository.save(savedUser);
		
		return AuthResponse.builder()
				.accessToken(accessToken)
				.refreshToken(refreshToken)
				.user(userMapper.toUserResponse(savedUser))
				.build();
	}
	
	@Override
	public AuthResponse login(LoginRequest request) {
		
		// Authenticate
		authenticationManager.authenticate(
				new UsernamePasswordAuthenticationToken(
						request.getEmail(),
						request.getPassword()
				)
		);
		
		User user = userRepository.findByEmail(request.getEmail())
				.orElseThrow(() -> new UserNotFoundException(
						"User not found: " + request.getEmail()
				));
		
		UserDetails userDetails = userDetailsService.loadUserByUsername(user.getEmail());
		String accessToken = jwtService.generateAccessToken(userDetails);
		String refreshToken = jwtService.generateRefreshToken(userDetails);
		
		// Update refresh token
		user.setRefreshToken(refreshToken);
		userRepository.save(user);
		
		return AuthResponse.builder()
				.accessToken(accessToken)
				.refreshToken(refreshToken)
				.user(userMapper.toUserResponse(user))
				.build();
	}
	
	@Override
	@Transactional
	public AuthResponse refreshToken(String refreshToken) {
		
		// Extract email from refresh token
		String email = jwtService.extractUsername(refreshToken);
		User user = userRepository.findByEmail(email)
				.orElseThrow(() -> new UserNotFoundException(
						"User not found: " + email
				));
		
		// Validate refresh token
		if (!refreshToken.equals(user.getRefreshToken())){
			throw new RuntimeException("Invalid refresh token");
		}
		
		UserDetails userDetails = userDetailsService.loadUserByUsername(email);
		if (!jwtService.isTokenValid(refreshToken, userDetails)){
			throw new RuntimeException("Refresh token expired");
		}
		
		String newAccessToken = jwtService.generateAccessToken(userDetails);
		String newRefreshToken = jwtService.generateRefreshToken(userDetails);
		
		return AuthResponse.builder()
				.accessToken(newAccessToken)
				.refreshToken(newRefreshToken)
				.user(userMapper.toUserResponse(user))
				.build();
	}
	
	@Override
	@Transactional
	public void logout(String email) {
		User user = userRepository.findByEmail(email)
				.orElseThrow(()-> new UserNotFoundException(
						"User not found: " + email
				));
		user.setRefreshToken(null);
		userRepository.save(user);
	}
}
