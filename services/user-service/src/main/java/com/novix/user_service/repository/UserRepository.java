package com.novix.user_service.repository;

import java.util.Optional;
import com.novix.user_service.entity.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface UserRepository extends JpaRepository<User, Long>{
	Optional<User> findByEmail(String email);
	Optional<User> findByUserName(String userName);
	Boolean existsByEmail(String email);
	Boolean existsByUserName(String userName);
}
