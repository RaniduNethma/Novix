package com.novix.user_service.config;

import java.util.Arrays;
import com.novix.user_service.entity.Role;
import com.novix.user_service.enums.RoleType;
import com.novix.user_service.repository.RoleRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.CommandLineRunner;
import org.springframework.stereotype.Component;

@Component
@Slf4j
@RequiredArgsConstructor
public class DataSeeder implements CommandLineRunner {
	private final RoleRepository roleRepository;
	
	@Override
	public void run(String... args) throws Exception {
		Arrays.stream(RoleType.values()).forEach(roleType -> {
			if (roleRepository.findByName(roleType).isEmpty()) {
				roleRepository.save(
						Role.builder()
								.name(roleType)
								.build()
				);
			}
		});
	}
}
