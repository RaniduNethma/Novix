package com.novix.user_service.repository;

import java.util.Optional;
import com.novix.user_service.entity.Role;
import com.novix.user_service.enums.RoleType;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface RoleRepository extends JpaRepository<Role, Long> {
	Optional<Role> findByName(RoleType name);
}
