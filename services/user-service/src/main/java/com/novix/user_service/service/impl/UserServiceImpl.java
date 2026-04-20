package com.novix.user_service.service.impl;

import com.novix.user_service.dto.request.UpdateProfileRequest;
import com.novix.user_service.dto.response.UserResponse;
import com.novix.user_service.entity.User;
import com.novix.user_service.exception.UserNotFoundException;
import com.novix.user_service.mapper.UserMapper;
import com.novix.user_service.repository.UserRepository;
import com.novix.user_service.service.UserService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@Slf4j
@RequiredArgsConstructor
public class UserServiceImpl implements UserService {
	private final UserRepository userRepository;
	private final UserMapper userMapper;
	
	@Override
	public UserResponse getUserById(Long id) {
		User user = userRepository.findById(id)
				.orElseThrow(() -> new UserNotFoundException(
						"User not found with id: " + id
				));
		return userMapper.toUserResponse(user);
	}
	
	@Override
	public UserResponse getUserByEmail(String email) {
		User user = userRepository.findByEmail(email)
				.orElseThrow(() -> new UserNotFoundException(
						"User not found with email: " + email
				));
		return userMapper.toUserResponse(user);
	}
	
	@Override
	@Transactional
	public UserResponse updateProfile(Long id, UpdateProfileRequest request) {
		User user = userRepository.findById(id)
				.orElseThrow(() -> new UserNotFoundException(
						"User not found with id: " + id
				));
		if (request.getFirstName() != null){
			user.setFirstName(request.getFirstName());
		}
		if (request.getLastName() != null){
			user.setLastName(request.getLastName());
		}
		if (request.getProfilePicture() != null){
			user.setProfilePicture(request.getProfilePicture());
		}
		User updatedUser = userRepository.save(user);
		return userMapper.toUserResponse(updatedUser);
	}
	
	@Override
	@Transactional
	public void deleteUser(Long id) {
		User user = userRepository.findById(id)
				.orElseThrow(() -> new UserNotFoundException(
						"User not found with id: " + id
				));
		userRepository.delete(user);
	}
}
