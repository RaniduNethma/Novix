package com.novix.user_service.mapper;

import java.util.Set;
import java.util.stream.Collectors;
import com.novix.user_service.dto.response.UserResponse;
import com.novix.user_service.entity.Role;
import com.novix.user_service.entity.User;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Named;

@Mapper(componentModel = "spring")
public interface UserMapper {
	
	@Mapping(target = "active", source = "active")
	@Mapping(target = "emailVerified", source = "emailVerified")
	@Mapping(target = "roles", source = "roles", qualifiedByName = "rolesToStrings")
	UserResponse toUserResponse(User user);
	
	@Named("rolesToStrings")
	default Set<String> rolesToStrings(Set<Role> roles) {
		if (roles == null) return Set.of();
		return roles.stream()
				.map(role -> role.getName().name())
				.collect(Collectors.toSet());
	}
}
